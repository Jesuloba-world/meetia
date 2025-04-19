"use client";

import { useState } from "react";
import { useRouter, useSearchParams } from "next/navigation";
import Link from "next/link";
import { useAuthStore } from "@/store/auth";
import { useMutation } from "@tanstack/react-query";

export default function Register() {
	const router = useRouter();
	const searchParams = useSearchParams();
	const redirect = searchParams.get("redirect") || "/";
	const { login } = useAuthStore();

	const [email, setEmail] = useState("");
	const [password, setPassword] = useState("");
	const [displayName, setDisplayName] = useState("");

	const registerMutation = useMutation({
		mutationFn: async () => {
			const response = await fetch(
				`${process.env.NEXT_PUBLIC_API_URL}/api/auth/register`,
				{
					method: "POST",
					headers: {
						"Content-Type": "application/json",
					},
					body: JSON.stringify({
						email,
						password,
						displayName,
					}),
				}
			);

			if (!response.ok) {
				if (response.status === 409) {
					throw new Error("User with this email already exists");
				} else {
					throw new Error("Registration failed");
				}
			}

			return response.json();
		},
		onSuccess: (data) => {
			login(data.user, data.token);
			router.push(redirect);
		},
	});

	const handleSubmit = (e: React.FormEvent) => {
		e.preventDefault();
		registerMutation.mutate();
	};

	return (
		<div className="min-h-screen flex items-center justify-center p-4">
			<div className="max-w-md w-full bg-white rounded-lg shadow p-8">
				<h1 className="text-2xl font-bold mb-6 text-center">
					Create Account
				</h1>

				<form onSubmit={handleSubmit}>
					<div className="mb-4">
						<label
							htmlFor="email"
							className="block text-sm font-medium text-gray-700 mb-1"
						>
							Email
						</label>
						<input
							type="email"
							id="email"
							value={email}
							onChange={(e) => setEmail(e.target.value)}
							className="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
							required
						/>
					</div>

					<div className="mb-4">
						<label
							htmlFor="displayName"
							className="block text-sm font-medium text-gray-700 mb-1"
						>
							Display Name
						</label>
						<input
							type="text"
							id="displayName"
							value={displayName}
							onChange={(e) => setDisplayName(e.target.value)}
							className="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
							required
						/>
					</div>

					<div className="mb-6">
						<label
							htmlFor="password"
							className="block text-sm font-medium text-gray-700 mb-1"
						>
							Password
						</label>
						<input
							type="password"
							id="password"
							value={password}
							onChange={(e) => setPassword(e.target.value)}
							className="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
							required
							minLength={6}
						/>
					</div>

					<button
						type="submit"
						className="w-full py-2 bg-blue-600 text-white rounded-md hover:bg-blue-700"
						disabled={registerMutation.isPending}
					>
						{registerMutation.isPending
							? "Creating Account..."
							: "Create Account"}
					</button>

					{registerMutation.isError && (
						<div className="mt-4 text-red-500 text-sm">
							{registerMutation.error.message ||
								"Registration failed"}
						</div>
					)}

					<div className="mt-4 text-center">
						<span className="text-sm text-gray-600">
							Already have an account?{" "}
							<Link
								href={`/login?redirect=${encodeURIComponent(
									redirect
								)}`}
								className="text-blue-600 hover:text-blue-800"
							>
								Sign In
							</Link>
						</span>
					</div>
				</form>
			</div>
		</div>
	);
}
