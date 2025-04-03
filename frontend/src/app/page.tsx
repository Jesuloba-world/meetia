import Link from "next/link";

export default function Home() {
	return (
		<main className="min-h-screen flex flex-col items-center justify-center p-4">
			<h1 className="text-4xl font-bold mb-8">Welcome to Meetia</h1>
			<div className="space-y-4 text-center">
				<p className="text-xl">Video conferencing made simple</p>
				<div className="flex gap-4 justify-center mt-6">
					<Link
						href="/new-meeting"
						className="px-4 py-2 bg-blue-600 text-white rounded-md hover:bg-blue-700 transition-colors"
					>
						New Meeting
					</Link>
					<a
						href="/join"
						className="px-4 py-2 border border-gray-300 rounded-md hover:bg-gray-100 transition-colors"
					>
						Join Meeting
					</a>
				</div>
			</div>
		</main>
	);
}
