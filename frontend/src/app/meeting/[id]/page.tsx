"use client";

import { useAuthStore } from "@/store/auth";
import { useParams, useRouter } from "next/navigation";
import { useEffect, useRef, useState } from "react";
import { cn } from "@/lib/utils";
import { useWebRTC } from "@/hooks/useWebRTC";
import {
	MicrophoneIcon,
	VideoCameraIcon,
	PhoneIcon,
	ComputerDesktopIcon,
	ChatBubbleLeftIcon,
} from "@heroicons/react/24/solid";
import {
	MicrophoneIcon as MicrophoneIconOutline,
	VideoCameraIcon as VideoCameraIconOutline,
} from "@heroicons/react/24/outline";

export default function MeetingRoom() {
	const params = useParams();
	const router = useRouter();
	const meetingId = params.id as string;
	const { user, isAuthenticated } = useAuthStore();
	const [showChat, setShowChat] = useState(false);

	const {
		isLoading,
		error,
		localStream,
		startLocalStream,
		isConnected,
		isVideoEnabled,
		isMuted,
		tracks,
		disconnect,
		toggleAudio,
		toggleVideo,
		isScreenSharing,
		screenStream,
		startScreenShare,
		stopScreenShare,
	} = useWebRTC(meetingId);

	const localVideoRef = useRef<HTMLVideoElement>(null);

	useEffect(() => {
		if (!isAuthenticated) {
			router.push(
				"/login?redirect=" + encodeURIComponent(`/meeting/${meetingId}`)
			);
		}
	}, [isAuthenticated, meetingId, router]);

	useEffect(() => {
		if (localStream && localVideoRef.current) {
			localVideoRef.current.srcObject = localStream;
		}
	}, [localStream]);

	useEffect(() => {
		if (isConnected && !localStream) {
			startLocalStream();
		}
	}, [isConnected, localStream, startLocalStream]);

	const handleEndCall = () => {
		disconnect();
		router.push("/");
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
					onClick={() => router.push("/")}
				>
					Back to Home
				</button>
			</div>
		);
	}

	return (
		<div className="flex flex-col h-screen">
			{/* Video Grid */}
			<div
				className={cn(
					"flex-1 grid gap-4 p-4 overflow-auto grid-cols-1 md:grid-cols-2",
					{ "grid-cols-1 md:grid-cols-2 lg:grid-cols-3": !showChat }
				)}
			>
				{/* Local video */}
				<div className="relative bg-gray-800 rounded-lg overflow-hidden">
					<video
						ref={localVideoRef}
						autoPlay
						playsInline
						muted
						className={cn("w-full h-full object-cover", {
							hidden: !isVideoEnabled,
						})}
					/>

					{!isVideoEnabled && (
						<div className="absolute inset-0 flex items-center justify-center bg-gray-900">
							<div className="w-24 h-24 rounded-full bg-gray-700 flex items-center justify-center">
								<span className="text-4xl">
									{user?.displayName?.[0] || "?"}
								</span>
							</div>
						</div>
					)}

					<div className="absolute bottom-2 left-2 text-white text-sm bg-black opacity-50 px-2 py-1 rounded">
						You {isMuted && "(Muted)"}
					</div>
				</div>

				{/* Remote participant videos */}
				{tracks
					.filter((track) => track.kind !== "audio") // only show video tracks
					.map((track) => (
						<div
							key={track.trackId}
							className="relative bg-gray-800 rounded-lg overflow-hidden"
						>
							<video
								autoPlay
								playsInline
								ref={(el) => {
									if (el) el.srcObject = track.stream;
								}}
								className="w-full h-full object-cover"
							/>
							<div className="absolute bottom-2 left-2 text-white text-sm bg-black opacity-50 px-2 py-1 rounded">
								{track.userId === user?.id
									? "You (Screen)"
									: `User ${track.userId}`}
							</div>
						</div>
					))}
			</div>

			{/* Chat Panel */}
			{showChat && (
				<div className="w-full md:w-1/3 border-l border-gray-300 h-full p-4">
					<div className="flex flex-col h-full">
						<div className="text-xl font-semibold mb-4">Chat</div>
						<div className="flex-1 overflow-y-auto mb-4 space-y-4">
							{/* TODO: complete the Chat part */}
							<div className="bg-gray-100 p-3 rounded-lg">
								<div className="font-semibold">John Doe</div>
								<div>Hello everyone!</div>
							</div>
						</div>
						<div className="mt-auto">
							<div className="flex">
								<input
									type="text"
									placeholder="Type a message..."
									className="flex-1 border border-gray-300 rounded-l-md px-3 py-2 focus:outline-none focus:ring-2 focus:ring-blue-500"
								/>
								<button className="bg-blue-600 text-white px-4 py-2 rounded-r-md">
									Send
								</button>
							</div>
						</div>
					</div>
				</div>
			)}

			{/* Control Bar */}
			<div className="bg-gray-100 p-4 flex items-center justify-center space-x-6">
				<button
					onClick={toggleAudio}
					className={`p-3 rounded-full ${
						isMuted
							? "bg-red-500 text-white"
							: "bg-white border border-gray-300"
					}`}
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
				>
					<ChatBubbleLeftIcon className="h-6 w-6" />
				</button>

				<button
					onClick={handleEndCall}
					className="p-3 rounded-full bg-red-500 text-white"
				>
					<PhoneIcon className="h-6 w-6 rotate-135" />
				</button>
			</div>
		</div>
	);
}
