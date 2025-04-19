import { useAuthStore } from "@/store/auth";
import { useCallback, useEffect, useRef, useState } from "react";
import useWebSocket, { ReadyState } from "react-use-websocket";

interface Track {
	stream: MediaStream;
	userId: string;
	kind: "audio" | "video" | "screen";
	trackId: string;
}

interface SignalMessage {
	type: "offer" | "answer" | "candidate";
	sdp?: string;
	candidate?: RTCIceCandidateInit;
	userId: string;
	meetingId: string;
	trackId?: string;
	target?: string;
}

interface ExtendedRTCPeerConnection extends RTCPeerConnection {
	_senders?: Map<string, RTCRtpSender>;
}

export function useWebRTC(meetingId: string) {
	const { user, token } = useAuthStore();
	const [isLoading, setIsLoading] = useState(true);
	const [error, setError] = useState<Error | null>(null);
	const [tracks, setTracks] = useState<Track[]>([]);
	const [localStream, setLocalStream] = useState<MediaStream | null>(null);
	const [screenStream, setScreenStream] = useState<MediaStream | null>(null);
	const [isMuted, setIsMuted] = useState(false);
	const [isVideoEnabled, setIsVideoEnabled] = useState(true);
	const [isScreenSharing, setIsScreenSharing] = useState(false);

	const peerConnection = useRef<ExtendedRTCPeerConnection | null>(null);
	const pendingSignalMessages = useRef<SignalMessage[]>([]);

	const serverUrl =
		process.env.NEXT_PUBLIC_API_URL || "http://localhost:8080";
	const wsUrl = `${serverUrl.replace(
		/^http/,
		"ws"
	)}/api/rtc/signal/${meetingId}`;

	const { sendMessage, lastMessage, readyState } = useWebSocket(
		token ? `${wsUrl}?token=${token}` : null,
		{
			shouldReconnect: () => true,
			reconnectAttempts: 10,
			reconnectInterval: 3000,
			onOpen: () => {
				console.log("websocket connection established");
				initializePeerConnection();
			},
			onError: (event) => {
				console.error("WebSocket error:", event);
				setError(new Error("WebSocket connection error"));
				setIsLoading(false);
			},
			onClose: (event) => {
				console.log("WebSocket connection closed", {
					code: event.code,
					reason: event.reason,
					wasClean: event.wasClean,
				});
				// If connection was closed abnormally, try to reconnect
				if (!event.wasClean) {
					console.log(
						"Connection closed abnormally, will attempt to reconnect"
					);
				}
			},
		}
	);

	const isConnected =
		readyState === ReadyState.OPEN && !!peerConnection.current;

	// send signal message through websocket
	const sendSignalMessage = useCallback(
		(message: SignalMessage) => {
			console.log(
				`Attempting to send ${message.type}, WebSocket state: ${readyState}`
			);
			if (readyState === ReadyState.OPEN) {
				console.log(`Actually sending: ${message.type}`);

				sendMessage(JSON.stringify(message));
			} else {
				console.warn(
					`WebSocket not open (state: ${readyState}) when trying to send ${message.type}. Message not sent.`
				);
				pendingSignalMessages.current.push(message); // queue message for later send
			}
		},
		[sendMessage, readyState]
	);

	const resendPendingSignals = useCallback(() => {
		if (pendingSignalMessages.current.length > 0) {
			if (readyState !== ReadyState.OPEN) {
				console.warn(
					`WebSocket not ready (state: ${readyState}), skipping resend`
				);
				return;
			}

			console.log(
				`Attempting to resend ${pendingSignalMessages.current.length} pending signals`
			);
			pendingSignalMessages.current.forEach((message) => {
				sendSignalMessage(message);
			});
			pendingSignalMessages.current = [];
		}
	}, [sendSignalMessage, readyState]);

	useEffect(() => {
		const interval = setInterval(() => {
			resendPendingSignals();
		}, 500);

		return () => clearInterval(interval);
	}, [resendPendingSignals]);

	const initializePeerConnection = useCallback(() => {
		if (!user) {
			console.log("Cannot initialize peer connection: No user");
			return;
		}

		pendingSignalMessages.current = [];

		const config: RTCConfiguration = {
			iceServers: [
				// {
				// 	urls: "turn:localhost:3478",
				// 	username: "meetia_user",
				// 	credential: "strong_password"
				// },
				{ urls: "stun:stun.l.google.com:19302" },
				{ urls: "stun:stun.l.google.com:5349" },
				{ urls: "stun:stun1.l.google.com:3478" },
				{ urls: "stun:stun1.l.google.com:5349" },
				{ urls: "stun:stun2.l.google.com:19302" },
				{ urls: "stun:stun2.l.google.com:5349" },
				{ urls: "stun:stun3.l.google.com:3478" },
				{ urls: "stun:stun3.l.google.com:5349" },
				{ urls: "stun:stun4.l.google.com:19302" },
				{ urls: "stun:stun4.l.google.com:5349" },
			],
			iceTransportPolicy: "all",
		};

		// create new peer connection
		const pc = new RTCPeerConnection(config);
		peerConnection.current = pc;

		pc.onconnectionstatechange = () => {
			console.log("Peer connection state changed:", pc.connectionState);
			if (
				pc.connectionState === "failed" ||
				pc.connectionState === "disconnected" ||
				pc.connectionState === "closed"
			) {
				console.error("Peer connection failed or disconnected");
			}
		};

		// Handle ICE candidate events
		pc.onicecandidate = (event) => {
			if (event.candidate) {
				if (
					["stable", "have-local-offer"].includes(pc.signalingState)
				) {
					console.log("Generated ICE candidate:", {
						type: event.candidate.type,
						protocol: event.candidate.protocol,
						address: event.candidate.address,
						port: event.candidate.port,
					});
					sendSignalMessage({
						type: "candidate",
						candidate: event.candidate.toJSON(),
						userId: user.id,
						meetingId,
					});
				} else {
					console.warn(
						"Dropping ICE candidate - invalid signaling state:",
						pc.signalingState
					);
				}
			}
		};

		pc.oniceconnectionstatechange = () => {
			console.log("ICE connection state:", pc.iceConnectionState);
		};

		// Handle track events
		pc.ontrack = (event) => {
			const stream = event.streams[0];
			const trackId = event.track.id;
			const userId = stream.id.split("-")[0];

			const kind =
				event.track.kind === "audio"
					? "audio"
					: stream.id.includes("screen")
					? "screen"
					: "video";

			const newTrack: Track = {
				stream,
				userId,
				kind,
				trackId,
			};

			setTracks((prev) => [
				...prev.filter((t) => t.trackId !== trackId),
				newTrack,
			]);

			setIsLoading(false);
		};
	}, [user, meetingId, sendSignalMessage]);

	const handleOffer = useCallback(
		async (message: SignalMessage) => {
			if (!peerConnection.current || !user || !message.sdp) return;

			try {
				await peerConnection.current.setRemoteDescription(
					new RTCSessionDescription({
						type: "offer",
						sdp: message.sdp,
					})
				);

				const answer = await peerConnection.current.createAnswer();
				await peerConnection.current.setLocalDescription(answer);

				sendSignalMessage({
					type: "answer",
					sdp: answer.sdp,
					userId: user.id,
					meetingId,
					target: message.userId,
				});
			} catch (err) {
				console.error("Error handling offer:", err);
			}
		},
		[meetingId, sendSignalMessage, user]
	);

	const handleAnswer = useCallback(
		async (message: SignalMessage) => {
			if (!peerConnection.current || !message.sdp) return;

			try {
				if (
					peerConnection.current.signalingState === "have-local-offer"
				) {
					await peerConnection.current.setRemoteDescription(
						new RTCSessionDescription({
							type: "answer",
							sdp: message.sdp,
						})
					);
					console.log(
						"Answer processed, new signaling state:",
						peerConnection.current.signalingState
					);
				} else {
					console.warn(
						"Received answer in unexpected state:",
						peerConnection.current.signalingState
					);

					// If connection is still active, ignore stale answer
					if (
						["connected", "connecting"].includes(
							peerConnection.current.connectionState
						)
					) {
						console.log(
							"Ignoring answer as connection is already established"
						);
						return;
					}

					console.log(
						"Restarting peer connection due to unexpected answer state"
					);
					initializePeerConnection();
				}
			} catch (err) {
				console.error("Error handling answer:", err);

				if (err instanceof Error) {
					initializePeerConnection();
				}
			}
		},
		[initializePeerConnection]
	);

	const handleCandidate = useCallback(async (message: SignalMessage) => {
		if (!peerConnection.current || !message.candidate) return;

		try {
			await peerConnection.current.addIceCandidate(
				new RTCIceCandidate(message.candidate)
			);
		} catch (err) {
			console.error("Error handling ICE candidate:", err);
		}
	}, []);

	// handle incoming signal message
	const handleSignalMessage = useCallback(
		(message: SignalMessage) => {
			if (!peerConnection.current) return;

			switch (message.type) {
				case "offer":
					handleOffer(message);
					break;
				case "answer":
					handleAnswer(message);
					break;
				case "candidate":
					handleCandidate(message);
					break;
			}
		},
		[handleOffer, handleAnswer, handleCandidate]
	);

	// process incoming websocket messages
	useEffect(() => {
		if (lastMessage && peerConnection.current) {
			try {
				const message = JSON.parse(lastMessage.data) as SignalMessage;
				console.log("Received WebSocket message:", message.type);
				handleSignalMessage(message);
				console.log(
					"Post-handler signaling state:",
					peerConnection.current.signalingState
				);
			} catch (err) {
				console.error("Error parsing WebSocket message:", err);
				console.error("Raw message data:", lastMessage.data);
			}
		}
	}, [lastMessage, handleSignalMessage]);

	// start local stream
	const startLocalStream = useCallback(
		async (video: boolean = true, audio: boolean = true) => {
			if (!peerConnection.current || !user) return null;

			try {
				const stream = await navigator.mediaDevices.getUserMedia({
					video,
					audio,
				});
				setLocalStream(stream);
				setIsMuted(!audio);
				setIsVideoEnabled(video);

				// Add tracks to peer connection
				stream.getTracks().forEach((track) => {
					if (peerConnection.current && stream) {
						const sender = peerConnection.current.addTrack(
							track,
							stream
						);

						// store sender for later reference
						if (!peerConnection.current._senders) {
							peerConnection.current._senders = new Map();
						}
						peerConnection.current._senders.set(track.kind, sender);
					}
				});

				// create and send offer
				const offer = await peerConnection.current.createOffer();
				await peerConnection.current.setLocalDescription(offer);

				sendSignalMessage({
					type: "offer",
					sdp: offer.sdp,
					userId: user.id,
					meetingId,
				});

				return stream;
			} catch (err) {
				console.error("Error accessing media devices:", err);
				setError(err instanceof Error ? err : new Error(String(err)));
				return null;
			}
		},
		[user, sendSignalMessage, meetingId]
	);

	// toggle audio mute state
	const toggleAudio = useCallback(() => {
		if (!localStream) return;

		const newMutedState = !isMuted;
		localStream.getAudioTracks().forEach((track) => {
			track.enabled = !newMutedState;
		});
		setIsMuted(newMutedState);
	}, [localStream, isMuted]);

	// toggle video
	const toggleVideo = useCallback(() => {
		if (!localStream) return;

		const newVideoState = !isVideoEnabled;
		localStream.getVideoTracks().forEach((track) => {
			track.enabled = newVideoState;
		});
		setIsVideoEnabled(newVideoState);
	}, [localStream, isVideoEnabled]);

	// stop screen sharing
	const stopScreenShare = useCallback(() => {
		if (!screenStream) return;

		screenStream.getTracks().forEach((track) => {
			track.stop();

			if (peerConnection.current) {
				// remove track from peer connection
				const senders = peerConnection.current.getSenders();
				const sender = senders.find((s) => s.track === track);
				if (sender) {
					peerConnection.current.removeTrack(sender);
				}
			}
		});
	}, [screenStream]);

	// start screem sharing
	const startScreenShare = useCallback(async () => {
		if (!peerConnection.current || !user) return null;

		try {
			const stream = await navigator.mediaDevices.getDisplayMedia({
				video: true,
			});
			setScreenStream(stream);
			setIsScreenSharing(true);

			const videoTrack = stream.getVideoTracks()[0];

			if (peerConnection.current && stream) {
				peerConnection.current.addTrack(videoTrack, stream);
			}

			// create and send offer
			const offer = await peerConnection.current.createOffer();
			await peerConnection.current.setLocalDescription(offer);

			sendSignalMessage({
				type: "offer",
				sdp: offer.sdp,
				userId: user.id,
				meetingId,
			});

			// handle track end
			videoTrack.onended = () => {
				stopScreenShare();
			};

			return stream;
		} catch (err) {
			console.error("Error starting screen share:", err);
			setError(err instanceof Error ? err : new Error(String(err)));
			return null;
		}
	}, [user, meetingId, sendSignalMessage, stopScreenShare]);

	// disconnect and cleanup
	const disconnect = useCallback(() => {
		// stop all media tracks
		if (localStream) {
			localStream.getTracks().forEach((track) => track.stop());
			setLocalStream(null);
		}

		if (screenStream) {
			localStream?.getTracks().forEach((track) => track.stop());
			setScreenStream(null);
		}

		// close peer connection
		if (peerConnection.current) {
			console.log("Closing peer connection in disconnect()");
			peerConnection.current.close();
			peerConnection.current = null;
		}

		pendingSignalMessages.current = [];
		setTracks([]);
	}, [localStream, screenStream]);

	useEffect(() => {
		console.log("useWebRTC hook mounted for meeting:", meetingId);
		return () => {
			console.log("useWebRTC hook unmounting for meeting:", meetingId);
		};
	}, [meetingId]);

	// cleanup on onmount
	useEffect(() => {
		return () => {
			disconnect();
		};
		// eslint-disable-next-line react-hooks/exhaustive-deps
	}, []);

	return {
		isConnected,
		isLoading,
		error,
		localStream,
		screenStream,
		tracks,
		isMuted,
		isVideoEnabled,
		isScreenSharing,
		startLocalStream,
		toggleAudio,
		toggleVideo,
		startScreenShare,
		stopScreenShare,
		disconnect,
	};
}
