"use client";

import { useAuthStore } from "@/store/auth";
import { useMutation } from "@tanstack/react-query";
import { useRouter } from "next/navigation";
import { useState } from "react";

export default function NewMeeting() {
	const router = useRouter();
	const { token } = useAuthStore();
	const [title, setTitle] = useState("");
	const [isPrivate, setIsPrivate] = useState(false);
	const [password, setPassword] = useState("");

	const createMeetingMutation = useMutation({
		mutationFn: async () => {
			const response = await fetch(
				`${process.env.NEXT_PUBLIC_API_URL}/api/meetings`,
				{
					method: "POST",
					headers: {
						"Content-Type": "application/json",
						Authorization: `Bearer ${token}`,
					},
					body: JSON.stringify({
						title,
						isPrivate,
						password: isPrivate ? password : undefined,
					}),
				}
			);

			if (!response.ok) {
				throw new Error("Failed to create meeting");
			}

			return response.json();
		},
		onSuccess: (data) => {
			router.push(`/meeting/${data.id}`);
		},
	});

	const handleSubmit = (e: React.FormEvent) => {
		e.preventDefault();
		createMeetingMutation.mutate();
	};

	return (
		<div className="min-h-screen flex items-center justify-center p-4">
			<div className="max-w-md w-full bg-white rounded-lg shadow p-8">
				<h1 className="text-2xl font-bold mb-6">Create New Meeting</h1>

				<form onSubmit={handleSubmit}>
					<div className="mb-4">
						<label
							htmlFor="title"
							className="block text-sm font-medium text-gray-700 mb-1"
						>
							Meeting Title
						</label>
						<input
							type="text"
							id="title"
							value={title}
							onChange={(e) => setTitle(e.target.value)}
							className="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
							required
						/>
					</div>

					<div className="mb-4">
						<div className="flex items-center">
							<input
								type="checkbox"
								id="isPrivate"
								checked={isPrivate}
								onChange={(e) => setIsPrivate(e.target.checked)}
								className="h-4 w-4 text-blue-600 focus:ring-blue-500 border-gray-300 rounded"
							/>
							<label
								htmlFor="isPrivate"
								className="ml-2 block text-sm text-gray-700"
							>
								Private Meeting (requires password)
							</label>
						</div>
					</div>

					{isPrivate && (
						<div className="mb-6">
							<label
								htmlFor="password"
								className="block text-sm font-medium text-gray-700 mb-1"
							>
								Meeting Password
							</label>
							<input
								type="password"
								id="password"
								value={password}
								onChange={(e) => setPassword(e.target.value)}
								className="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
								required={isPrivate}
							/>
						</div>
					)}

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
							disabled={createMeetingMutation.isPending}
						>
							{createMeetingMutation.isPending
								? "Creating..."
								: "Create Meeting"}
						</button>
					</div>

					{createMeetingMutation.isError && (
						<div className="mt-4 text-red-500 text-sm">
							{createMeetingMutation.error.message ||
								"Failed to create meeting"}
						</div>
					)}
				</form>
			</div>
		</div>
	);
}
