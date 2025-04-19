import { useState, useRef, useEffect } from "react";
import { useAuthStore } from "@/store/auth";
import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { PaperAirplaneIcon } from "@heroicons/react/24/solid";

interface ChatMessage {
	id: string;
	meetingId: string;
	userId: string;
	message: string;
	sentAt: string;
	user: {
		displayName: string;
	};
}

interface MeetingChatProps {
	meetingId: string;
}

export default function MeetingChat({ meetingId }: MeetingChatProps) {
	const { user, token } = useAuthStore();
	const [message, setMessage] = useState("");
	const messagesEndRef = useRef<HTMLDivElement>(null);
	const queryClient = useQueryClient();

	// Auto-scroll to bottom when new messages arrive
	const scrollToBottom = () => {
		messagesEndRef.current?.scrollIntoView({ behavior: "smooth" });
	};

	// Fetch chat messages
	const { data: chatMessages, isLoading } = useQuery<ChatMessage[]>({
		queryKey: ["chat", meetingId],
		queryFn: async () => {
			const response = await fetch(
				`${process.env.NEXT_PUBLIC_API_URL}/api/meetings/${meetingId}/chat`,
				{
					headers: {
						Authorization: `Bearer ${token}`,
					},
				}
			);

			if (!response.ok) {
				throw new Error("Failed to fetch chat messages");
			}

			return response.json();
		},
		refetchInterval: 5000, // Poll for new messages every 5 seconds
	});

	// Send chat message
	const sendMessageMutation = useMutation({
		mutationFn: async (text: string) => {
			const response = await fetch(
				`${process.env.NEXT_PUBLIC_API_URL}/api/meetings/${meetingId}/chat`,
				{
					method: "POST",
					headers: {
						"Content-Type": "application/json",
						Authorization: `Bearer ${token}`,
					},
					body: JSON.stringify({
						message: text,
					}),
				}
			);

			if (!response.ok) {
				throw new Error("Failed to send message");
			}

			return response;
		},
		onSuccess: () => {
			// Clear input and refetch messages
			setMessage("");
			queryClient.invalidateQueries({ queryKey: ["chat", meetingId] });
		},
	});

	// Handle form submission
	const handleSubmit = (e: React.FormEvent) => {
		e.preventDefault();
		if (message.trim()) {
			sendMessageMutation.mutate(message);
		}
	};

	// Scroll to bottom when messages change
	useEffect(() => {
		scrollToBottom();
	}, [chatMessages]);

	// Format timestamp
	const formatTime = (timestamp: string) => {
		const date = new Date(timestamp);
		return date.toLocaleTimeString([], {
			hour: "2-digit",
			minute: "2-digit",
		});
	};

	return (
		<div className="flex flex-col h-full">
			<div className="text-xl font-semibold mb-4">Chat</div>

			<div className="flex-1 overflow-y-auto mb-4 space-y-4">
				{isLoading ? (
					<div className="text-center text-gray-500">
						Loading messages...
					</div>
				) : chatMessages && chatMessages.length > 0 ? (
					chatMessages.map((msg) => (
						<div
							key={msg.id}
							className={`p-3 rounded-lg max-w-[80%] ${
								msg.userId === user?.id
									? "ml-auto bg-blue-500 text-white"
									: "bg-gray-100 text-gray-800"
							}`}
						>
							<div className="font-semibold text-sm">
								{msg.userId === user?.id
									? "You"
									: msg.user.displayName}
							</div>
							<div>{msg.message}</div>
							<div className="text-xs opacity-70 text-right mt-1">
								{formatTime(msg.sentAt)}
							</div>
						</div>
					))
				) : (
					<div className="text-center text-gray-500">
						No messages yet
					</div>
				)}
				<div ref={messagesEndRef} />
			</div>

			<form onSubmit={handleSubmit} className="mt-auto">
				<div className="flex">
					<input
						type="text"
						value={message}
						onChange={(e) => setMessage(e.target.value)}
						placeholder="Type a message..."
						className="flex-1 border border-gray-300 rounded-l-md px-3 py-2 focus:outline-none focus:ring-2 focus:ring-blue-500"
					/>
					<button
						type="submit"
						className="bg-blue-600 text-white px-4 py-2 rounded-r-md hover:bg-blue-700 transition-colors flex items-center"
						disabled={sendMessageMutation.isPending}
					>
						<PaperAirplaneIcon className="h-5 w-5" />
					</button>
				</div>
				{sendMessageMutation.isError && (
					<div className="text-red-500 text-sm mt-1">
						Failed to send message. Please try again.
					</div>
				)}
			</form>
		</div>
	);
}
