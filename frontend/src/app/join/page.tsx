"use client";

import { useState } from "react";
import { useRouter, useSearchParams } from "next/navigation";
import { useAuthStore } from "@/store/auth";
import { useMutation } from "@tanstack/react-query";

export default function JoinMeeting() {
	const router = useRouter();
	const searchParams = useSearchParams();
	const { token, isAuthenticated } = useAuthStore();

	const initialCode = searchParams.get("code") || "";
	const [meetingCode, setMeetingCode] = useState(initialCode);
	const [password, setPassword] = useState("");

	const joinMeetingMutation = useMutation({
		mutationFn: async () => {
			const response = await fetch(
				`${process.env.NEXT_PUBLIC_API_URL}/api/meetings/join`,
				{
					method: "POST",
					headers: {
						"Content-Type": "application/json",
						Authorization: `Bearer ${token}`,
					},
					body: JSON.stringify({
						meetingCode,
						password,
					}),
				}
			);

			if (!response.ok) {
				if (response.status === 401) {
					throw new Error("Invalid password");
				} else if (response.status === 404) {
					throw new Error("Meeting not found");
				} else {
					throw new Error("Failed to join meeting");
				}
			}

			return response.json();
		},
		onSuccess: (data) => {
			router.push(`/meeting/${data.id}`);
		},
	});

	// Redirect to login if not authenticated
	if (!isAuthenticated) {
		router.push(
			"/login?redirect=" +
				encodeURIComponent(
					"/join" + (initialCode ? `?code=${initialCode}` : "")
				)
		);
		return null;
	}

	const handleSubmit = (e: React.FormEvent) => {
		e.preventDefault();
		joinMeetingMutation.mutate();
	};

	return (
		<div className="min-h-screen flex items-center justify-center p-4">
			<div className="max-w-md w-full bg-white rounded-lg shadow p-8">
				<h1 className="text-2xl font-bold mb-6">Join Meeting</h1>

				<form onSubmit={handleSubmit}>
					<div className="mb-4">
						<label
							htmlFor="meetingCode"
							className="block text-sm font-medium text-gray-700 mb-1"
						>
							Meeting Code
						</label>
						<input
							type="text"
							id="meetingCode"
							value={meetingCode}
							onChange={(e) => setMeetingCode(e.target.value)}
							className="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
							placeholder="Enter meeting code"
							required
						/>
					</div>

					<div className="mb-6">
						<label
							htmlFor="password"
							className="block text-sm font-medium text-gray-700 mb-1"
						>
							Meeting Password (if required)
						</label>
						<input
							type="password"
							id="password"
							value={password}
							onChange={(e) => setPassword(e.target.value)}
							className="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
							placeholder="Enter password if meeting is private"
						/>
					</div>

					<div className="flex items-center justify-between">
						<button
							type="button"
							onClick={() => router.push("/")}
							className="px-4 py-2 border border-gray-300 rounded-md hover:bg-gray-50"
						>
							Cancel
						</button>
						<button
							type="submit"
							className="px-4 py-2 bg-blue-600 text-white rounded-md hover:bg-blue-700"
							disabled={joinMeetingMutation.isPending}
						>
							{joinMeetingMutation.isPending
								? "Joining..."
								: "Join Meeting"}
						</button>
					</div>

					{joinMeetingMutation.isError && (
						<div className="mt-4 text-red-500 text-sm">
							{joinMeetingMutation.error.message ||
								"Failed to join meeting"}
						</div>
					)}
				</form>
			</div>
		</div>
	);
}
