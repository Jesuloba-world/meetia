export default function LoadingSpinner({
	size = "md",
}: {
	size?: "sm" | "md" | "lg";
}) {
	const sizeClasses = {
		sm: "h-4 w-4 border-2",
		md: "h-8 w-8 border-4",
		lg: "h-12 w-12 border-[6px]",
	};

	return (
		<div className="flex justify-center items-center">
			<div
				className={`${sizeClasses[size]} rounded-full border-blue-600 border-t-transparent animate-spin`}
			></div>
		</div>
	);
}
