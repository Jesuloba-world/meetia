// frontend/src/app/page.tsx
import Link from "next/link";
import {
	VideoCameraIcon,
	UserGroupIcon,
	ChartBarIcon,
	LockClosedIcon,
} from "@heroicons/react/24/outline";

export default function Home() {
	return (
		<div className="min-h-screen">
			{/* Hero Section */}
			<div className="bg-gradient-to-r from-blue-600 to-indigo-700 text-white">
				<div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-24">
					<div className="text-center">
						<h1 className="text-4xl md:text-5xl font-bold mb-6">
							Video conferencing made simple
						</h1>
						<p className="text-xl mb-10 max-w-3xl mx-auto">
							Secure, reliable, and high-quality video meetings
							for teams of all sizes.
						</p>
						<div className="flex flex-col sm:flex-row gap-4 justify-center">
							<Link
								href="/new-meeting"
								className="px-6 py-3 bg-white text-blue-700 font-medium rounded-md hover:bg-gray-100 transition-colors"
							>
								New Meeting
							</Link>
							<Link
								href="/join"
								className="px-6 py-3 border border-white text-white font-medium rounded-md hover:bg-white/10 transition-colors"
							>
								Join Meeting
							</Link>
						</div>
					</div>
				</div>
			</div>

			{/* Features Section */}
			<div className="py-16 bg-white">
				<div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
					<div className="text-center mb-16">
						<h2 className="text-3xl font-bold text-gray-900">
							Why choose Meetia?
						</h2>
						<p className="mt-4 text-lg text-gray-600 max-w-3xl mx-auto">
							Our platform offers everything you need for
							successful video meetings.
						</p>
					</div>

					<div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-8">
						<div className="p-6 bg-gray-50 rounded-lg">
							<VideoCameraIcon className="h-12 w-12 text-blue-600 mb-4" />
							<h3 className="text-xl font-semibold mb-2">
								HD Video Quality
							</h3>
							<p className="text-gray-600">
								Crystal clear video even with low bandwidth
								connections.
							</p>
						</div>

						<div className="p-6 bg-gray-50 rounded-lg">
							<UserGroupIcon className="h-12 w-12 text-blue-600 mb-4" />
							<h3 className="text-xl font-semibold mb-2">
								Multiple Participants
							</h3>
							<p className="text-gray-600">
								Host meetings with multiple participants without
								quality issues.
							</p>
						</div>

						<div className="p-6 bg-gray-50 rounded-lg">
							<ChartBarIcon className="h-12 w-12 text-blue-600 mb-4" />
							<h3 className="text-xl font-semibold mb-2">
								Screen Sharing
							</h3>
							<p className="text-gray-600">
								Share your screen with participants for
								effective presentations.
							</p>
						</div>

						<div className="p-6 bg-gray-50 rounded-lg">
							<LockClosedIcon className="h-12 w-12 text-blue-600 mb-4" />
							<h3 className="text-xl font-semibold mb-2">
								Secure Meetings
							</h3>
							<p className="text-gray-600">
								Password protection and encrypted connections
								for all meetings.
							</p>
						</div>
					</div>
				</div>
			</div>

			{/* CTA Section */}
			<div className="bg-gray-50 py-16">
				<div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 text-center">
					<h2 className="text-3xl font-bold text-gray-900 mb-4">
						Ready to start your meeting?
					</h2>
					<p className="text-xl text-gray-600 mb-8 max-w-3xl mx-auto">
						Get started with Meetia today and experience better
						video meetings.
					</p>
					<div className="flex flex-col sm:flex-row gap-4 justify-center">
						<Link
							href="/register"
							className="px-6 py-3 bg-blue-600 text-white font-medium rounded-md hover:bg-blue-700 transition-colors"
						>
							Sign Up Free
						</Link>
						<Link
							href="/login"
							className="px-6 py-3 border border-gray-300 text-gray-700 font-medium rounded-md hover:bg-gray-100 transition-colors"
						>
							Sign In
						</Link>
					</div>
				</div>
			</div>

			{/* Footer */}
			<footer className="bg-gray-800 text-white py-12">
				<div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
					<div className="flex flex-col md:flex-row justify-between items-center">
						<div className="mb-4 md:mb-0">
							<div className="flex items-center">
								<VideoCameraIcon className="h-8 w-8 text-blue-400" />
								<span className="ml-2 text-xl font-bold">
									Meetia
								</span>
							</div>
							<p className="mt-2 text-gray-400">
								Video conferencing for everyone
							</p>
						</div>
						<div className="flex space-x-6">
							<a
								href="#"
								className="text-gray-400 hover:text-white"
							>
								Privacy
							</a>
							<a
								href="#"
								className="text-gray-400 hover:text-white"
							>
								Terms
							</a>
							<a
								href="#"
								className="text-gray-400 hover:text-white"
							>
								Help
							</a>
						</div>
					</div>
					<div className="mt-8 text-center text-gray-400 text-sm">
						&copy; {new Date().getFullYear()} Meetia. All rights
						reserved.
					</div>
				</div>
			</footer>
		</div>
	);
}
