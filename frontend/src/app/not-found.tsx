import Link from "next/link";

export default function NotFound() {
	return (
		<div className="min-h-screen flex items-center justify-center p-4">
			<div className="text-center">
				<h1 className="text-6xl font-bold text-gray-900 mb-4">404</h1>
				<h2 className="text-2xl font-semibold text-gray-700 mb-6">
					Page Not Found
				</h2>
				<p className="text-gray-600 mb-8 max-w-md mx-auto">
					The page you&apos;re looking for doesn&apos;t exist or has
					been moved.
				</p>
				<Link
					href="/"
					className="px-6 py-3 bg-blue-600 text-white font-medium rounded-md hover:bg-blue-700 transition-colors"
				>
					Back to Home
				</Link>
			</div>
		</div>
	);
}
