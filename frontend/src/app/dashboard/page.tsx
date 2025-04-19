"use client";

import { useEffect } from "react";
import { useRouter } from "next/navigation";
import Link from "next/link";
import { useAuthStore } from "@/store/auth";
import { useQuery } from "@tanstack/react-query";
import {
	CalendarIcon,
	VideoCameraIcon,
	UserGroupIcon,
	ClockIcon,
} from "@heroicons/react/24/outline";

interface Meeting {
	id: string;
	title: string;
	meetingCode: string;
	isPrivate: boolean;
	createdAt: string;
	hostId: string;
	host: {
		displayName: string;
	};
}

export default function Dashboard() {
	const router = useRouter();
	const { user, token, isAuthenticated } = useAuthStore();

	useEffect(() => {
		if (!isAuthenticated) {
			router.push("/login?redirect=/dashboard");
		}
	}, [isAuthenticated, router]);

	const {
		data: meetings,
		isLoading,
		error,
	} = useQuery<Meeting[]>({
		queryKey: ["meetings"],
		queryFn: async () => {
			const response = await fetch(
				`${process.env.NEXT_PUBLIC_API_URL}/api/meetings`,
				{
					headers: {
						Authorization: `Bearer ${token}`,
					},
				}
			);

			if (!response.ok) {
				throw new Error("Failed to fetch meetings");
			}

			return response.json();
		},
		enabled: !!token,
	});

	if (!isAuthenticated) {
		return null;
	}

	return (
		<div className="min-h-screen bg-gray-50">
			<header className="bg-white shadow">
				<div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-6">
					<div className="flex justify-between items-center">
						<h1 className="text-2xl font-bold text-gray-900">
							Your Meetings
						</h1>
						<div className="flex space-x-4">
							<Link
								href="/new-meeting"
								className="px-4 py-2 bg-blue-600 text-white rounded-md hover:bg-blue-700 transition-colors"
							>
								New Meeting
							</Link>
							<Link
								href="/join"
								className="px-4 py-2 border border-gray-300 rounded-md hover:bg-gray-100 transition-colors"
							>
								Join Meeting
							</Link>
						</div>
					</div>
				</div>
			</header>

			<main className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
				{isLoading ? (
					<div className="text-center py-10">
						<p className="text-gray-600">
							Loading your meetings...
						</p>
					</div>
				) : error ? (
					<div className="text-center py-10">
						<p className="text-red-600">Failed to load meetings.</p>
					</div>
				) : meetings && meetings.length > 0 ? (
					<div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
						{meetings.map((meeting) => (
							<div
								key={meeting.id}
								className="bg-white shadow rounded-lg overflow-hidden"
							>
								<div className="p-6">
									<h2 className="text-xl font-semibold mb-2">
										{meeting.title}
									</h2>

									<div className="space-y-2 mb-4">
										<div className="flex items-center text-sm text-gray-600">
											<UserGroupIcon className="h-5 w-5 mr-2" />
											<span>
												Hosted by{" "}
												{meeting.hostId === user?.id
													? "you"
													: meeting.host.displayName}
											</span>
										</div>

										<div className="flex items-center text-sm text-gray-600">
											<ClockIcon className="h-5 w-5 mr-2" />
											<span>
												Created{" "}
												{new Date(
													meeting.createdAt
												).toLocaleString()}
											</span>
										</div>

										<div className="flex items-center text-sm text-gray-600">
											<CalendarIcon className="h-5 w-5 mr-2" />
											<span>
												Meeting code:{" "}
												{meeting.meetingCode}
											</span>
										</div>

										{meeting.isPrivate && (
											<div className="text-sm text-amber-600">
												Password protected
											</div>
										)}
									</div>

									<div className="flex justify-end">
										<Link
											href={`/meeting/${meeting.id}`}
											className="flex items-center px-4 py-2 bg-blue-600 text-white rounded-md hover:bg-blue-700 transition-colors"
										>
											<VideoCameraIcon className="h-5 w-5 mr-2" />
											Join
										</Link>
									</div>
								</div>
							</div>
						))}
					</div>
				) : (
					<div className="text-center py-16 bg-white rounded-lg shadow">
						<h3 className="text-lg font-medium text-gray-900 mb-1">
							No meetings yet
						</h3>
						<p className="text-gray-600 mb-6">
							Create a new meeting or join an existing one.
						</p>
						<div className="flex justify-center space-x-4">
							<Link
								href="/new-meeting"
								className="px-4 py-2 bg-blue-600 text-white rounded-md hover:bg-blue-700 transition-colors"
							>
								New Meeting
							</Link>
							<Link
								href="/join"
								className="px-4 py-2 border border-gray-300 rounded-md hover:bg-gray-100 transition-colors"
							>
								Join Meeting
							</Link>
						</div>
					</div>
				)}
			</main>
		</div>
	);
}
