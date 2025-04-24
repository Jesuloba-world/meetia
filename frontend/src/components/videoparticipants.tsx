import { useRef, useEffect } from "react";

interface VideoParticipantProps {
	stream: MediaStream;
	displayName: string;
	isMuted?: boolean;
	isLocal?: boolean;
	isSpeaking?: boolean;
}

export default function VideoParticipant({
	stream,
	displayName,
	isMuted = false,
	isLocal = false,
	isSpeaking = false,
}: VideoParticipantProps) {
	const videoRef = useRef<HTMLVideoElement>(null);

	useEffect(() => {
		if (videoRef.current && stream) {
			videoRef.current.srcObject = stream;
		}
	}, [stream]);

	// Check if stream has video tracks
	const hasVideo =
		stream.getVideoTracks().length > 0 &&
		stream.getVideoTracks()[0].enabled;

	return (
		<div
			className={`relative bg-gray-800 rounded-lg overflow-hidden ${
				isSpeaking ? "ring-4 ring-green-500" : ""
			}`}
		>
			{hasVideo ? (
				<video
					ref={videoRef}
					autoPlay
					playsInline
					muted={isLocal || isMuted}
					className={`w-full h-full object-cover ${
						isLocal ? "scale-x-[-1]" : ""
					}`}
				/>
			) : (
				<div className="h-full flex items-center justify-center bg-gray-900 min-h-[200px]">
					<div className="w-24 h-24 rounded-full bg-gray-700 flex items-center justify-center">
						<span className="text-4xl text-white">
							{displayName[0] || "?"}
						</span>
					</div>
				</div>
			)}

			<div className="absolute bottom-2 left-2 right-2 flex justify-between items-center">
				<div className="text-white text-sm bg-black bg-opacity-50 px-2 py-1 rounded">
					{displayName} {isLocal ? "(You)" : ""}
				</div>

				{isMuted && (
					<div className="bg-red-500 text-white text-xs px-2 py-1 rounded-full">
						Muted
					</div>
				)}
			</div>
		</div>
	);
}
