"use client";

import { useEffect, useState } from "react";
import { useParams, useRouter } from "next/navigation";
import { useWebRTC } from "@/hooks/useWebRTC";
import { useAuthStore } from "@/store/auth";
import {
	MicrophoneIcon,
	VideoCameraIcon,
	PhoneIcon,
	ComputerDesktopIcon,
	ChatBubbleLeftIcon,
	XMarkIcon,
} from "@heroicons/react/24/solid";
import {
	MicrophoneIcon as MicrophoneIconOutline,
	VideoCameraIcon as VideoCameraIconOutline,
} from "@heroicons/react/24/outline";
import VideoParticipant from "@/components/videoparticipants";
import MeetingChat from "@/components/meetingChat";
import { useQuery } from "@tanstack/react-query";

interface Meeting {
	id: string;
	title: string;
	meetingCode: string;
	host: {
		id: string;
		displayName: string;
	};
}

interface Participant {
	id: string;
	userId: string;
	user: {
		displayName: string;
	};
	role: string;
}

export default function MeetingRoom() {
	const params = useParams();
	const router = useRouter();
	const meetingId = params.id as string;
	const { user, token, isAuthenticated, _hasHydrated } = useAuthStore();
	const [showChat, setShowChat] = useState(false);

	// Fetch meeting details
	const { data: meeting } = useQuery<Meeting>({
		queryKey: ["meeting", meetingId],
		queryFn: async () => {
			const response = await fetch(
				`${process.env.NEXT_PUBLIC_API_URL}/api/meetings/${meetingId}`,
				{
					headers: {
						Authorization: `Bearer ${token}`,
					},
				}
			);

			if (!response.ok) {
				throw new Error("Failed to fetch meeting details");
			}

			return response.json().then((data) => data.meeting);
		},
		enabled: !!token,
	});

	// Fetch participants
	const { data: participants } = useQuery<Participant[]>({
		queryKey: ["participants", meetingId],
		queryFn: async () => {
			const response = await fetch(
				`${process.env.NEXT_PUBLIC_API_URL}/api/meetings/${meetingId}/participants`,
				{
					headers: {
						Authorization: `Bearer ${token}`,
					},
				}
			);

			if (!response.ok) {
				throw new Error("Failed to fetch participants");
			}

			return response.json().then((data) => {
				return data.participants;
			});
		},
		enabled: !!token,
		refetchInterval: 10000, // Poll for new participants every 10 seconds
	});

	console.log(participants);

	const {
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
	} = useWebRTC(meetingId);

	useEffect(() => {
		if (!isAuthenticated && _hasHydrated) {
			router.push(
				"/login?redirect=" + encodeURIComponent(`/meeting/${meetingId}`)
			);
		}
	}, [isAuthenticated, router, meetingId, _hasHydrated]);

	useEffect(() => {
		if (isConnected && !localStream) {
			startLocalStream();
		}
	}, [isConnected, localStream, startLocalStream]);

	const handleEndCall = () => {
		disconnect();
		router.push("/dashboard");
	};

	// Helper function to get display name for a track
	const getDisplayName = (userId: string) => {
		if (userId === user?.id) return user.displayName;

		const participant = participants?.find((p) => p.userId === userId);
		return participant?.user.displayName || "Unknown User";
	};

	if (isLoading) {
		return (
			<div className="flex items-center justify-center min-h-screen">
				<div className="text-lg">Connecting to meeting...</div>
			</div>
		);
	}

	if (error) {
		return (
			<div className="flex flex-col items-center justify-center min-h-screen">
				<div className="text-red-500 text-xl mb-4">
					Error connecting to meeting
				</div>
				<div className="mb-8">{error.message}</div>
				<button
					className="px-4 py-2 bg-blue-600 text-white rounded-md"
					onClick={() => router.push("/dashboard")}
				>
					Back to Dashboard
				</button>
			</div>
		);
	}

	return (
		<div className="flex flex-col h-screen">
			{/* Meeting Header */}
			<div className="bg-white shadow py-3 px-6">
				<div className="flex items-center justify-between">
					<div>
						<h1 className="font-semibold">
							{meeting?.title || "Meeting"}
						</h1>
						<div className="text-sm text-gray-500">
							Meeting code: {meeting?.meetingCode}
						</div>
					</div>
					<div className="flex items-center">
						<div className="text-sm mr-4">
							{new Date().toLocaleTimeString([], {
								hour: "2-digit",
								minute: "2-digit",
							})}
						</div>
					</div>
				</div>
			</div>

			{/* Main Content Area */}
			<div className="flex flex-1 overflow-hidden">
				{/* Video Grid */}
				<div
					className={`flex-1 p-4 overflow-auto ${
						showChat ? "sm:pr-0" : ""
					}`}
				>
					<div
						className={`grid gap-4 h-full 
            ${
				tracks.length === 0
					? "grid-cols-1"
					: tracks.length < 2
					? "grid-cols-1 md:grid-cols-2"
					: "grid-cols-1 md:grid-cols-2 lg:grid-cols-3"
			}`}
					>
						{/* Local Video */}
						{localStream && (
							<VideoParticipant
								stream={localStream}
								displayName={user?.displayName || "You"}
								isMuted={isMuted}
								isLocal={true}
							/>
						)}

						{/* Screen Share */}
						{screenStream && (
							<VideoParticipant
								stream={screenStream}
								displayName={`${
									user?.displayName || "Your"
								} Screen`}
								isLocal={true}
							/>
						)}

						{/* Remote Participant Videos */}
						{tracks
							.filter((track) => track.kind !== "audio") // Only show video tracks
							.map((track) => (
								<VideoParticipant
									key={track.trackId}
									stream={track.stream}
									displayName={getDisplayName(track.userId)}
									isMuted={false} // We don't know if they're muted yet
								/>
							))}
					</div>
				</div>

				{/* Chat Panel */}
				{showChat && (
					<div className="hidden sm:block sm:w-80 md:w-96 border-l border-gray-300 h-full p-4 bg-white">
						<MeetingChat meetingId={meetingId} />
					</div>
				)}

				{/* Mobile Chat Overlay */}
				{showChat && (
					<div className="sm:hidden fixed inset-0 bg-white z-10 p-4">
						<div className="h-full flex flex-col">
							<div className="flex justify-between items-center mb-4">
								<h2 className="font-semibold text-lg">Chat</h2>
								<button
									className="p-2 rounded-full hover:bg-gray-100"
									onClick={() => setShowChat(false)}
								>
									<XMarkIcon className="h-6 w-6" />
								</button>
							</div>
							<MeetingChat meetingId={meetingId} />
						</div>
					</div>
				)}
			</div>

			{/* Control Bar */}
			<div className="bg-gray-100 p-4 flex items-center justify-center space-x-6">
				<button
					onClick={toggleAudio}
					className={`p-3 rounded-full ${
						isMuted
							? "bg-red-500 text-white"
							: "bg-white border border-gray-300"
					}`}
					aria-label={isMuted ? "Unmute" : "Mute"}
				>
					{isMuted ? (
						<MicrophoneIconOutline className="h-6 w-6" />
					) : (
						<MicrophoneIcon className="h-6 w-6" />
					)}
				</button>

				<button
					onClick={toggleVideo}
					className={`p-3 rounded-full ${
						!isVideoEnabled
							? "bg-red-500 text-white"
							: "bg-white border border-gray-300"
					}`}
					aria-label={
						isVideoEnabled ? "Turn off camera" : "Turn on camera"
					}
				>
					{!isVideoEnabled ? (
						<VideoCameraIconOutline className="h-6 w-6" />
					) : (
						<VideoCameraIcon className="h-6 w-6" />
					)}
				</button>

				<button
					onClick={
						isScreenSharing ? stopScreenShare : startScreenShare
					}
					className={`p-3 rounded-full ${
						isScreenSharing
							? "bg-blue-500 text-white"
							: "bg-white border border-gray-300"
					}`}
					aria-label={
						isScreenSharing ? "Stop sharing screen" : "Share screen"
					}
				>
					<ComputerDesktopIcon className="h-6 w-6" />
				</button>

				<button
					onClick={() => setShowChat(!showChat)}
					className={`p-3 rounded-full ${
						showChat
							? "bg-blue-500 text-white"
							: "bg-white border border-gray-300"
					}`}
					aria-label={showChat ? "Hide chat" : "Show chat"}
				>
					<ChatBubbleLeftIcon className="h-6 w-6" />
				</button>

				<button
					onClick={handleEndCall}
					className="p-3 rounded-full bg-red-500 text-white"
					aria-label="End call"
				>
					<PhoneIcon className="h-6 w-6 transform rotate-135" />
				</button>
			</div>
		</div>
	);
}
