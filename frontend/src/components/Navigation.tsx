// frontend/src/components/Navigation.tsx
"use client";

import { useState } from "react";
import Link from "next/link";
import { usePathname } from "next/navigation";
import { useAuthStore } from "@/store/auth";
import {
	VideoCameraIcon,
	Bars3Icon,
	XMarkIcon,
} from "@heroicons/react/24/outline";

export default function Navigation() {
	const pathname = usePathname();
	const { user, isAuthenticated, logout } = useAuthStore();
	const [isMobileMenuOpen, setIsMobileMenuOpen] = useState(false);

	const navItems = [
		{ name: "Home", href: "/" },
		{ name: "Dashboard", href: "/dashboard" },
	];

	if (
		pathname === "/login" ||
		pathname === "/register" ||
		pathname.startsWith("/meeting/")
	) {
		return null;
	}

	return (
		<nav className="bg-white shadow">
			<div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
				<div className="flex justify-between h-16">
					<div className="flex items-center">
						<Link
							href="/"
							className="flex-shrink-0 flex items-center"
						>
							<VideoCameraIcon className="h-8 w-8 text-blue-600" />
							<span className="ml-2 text-xl font-bold text-gray-900">
								Meetia
							</span>
						</Link>

						<div className="hidden sm:ml-6 sm:flex sm:space-x-8">
							{navItems.map((item) => (
								<Link
									key={item.href}
									href={item.href}
									className={`inline-flex items-center px-1 pt-1 border-b-2 text-sm font-medium ${
										pathname === item.href
											? "border-blue-500 text-gray-900"
											: "border-transparent text-gray-500 hover:border-gray-300 hover:text-gray-700"
									}`}
								>
									{item.name}
								</Link>
							))}
						</div>
					</div>

					<div className="hidden sm:ml-6 sm:flex sm:items-center">
						{isAuthenticated ? (
							<div className="flex items-center">
								<span className="px-3 py-2 text-sm text-gray-700">
									Hi, {user?.displayName}
								</span>
								<button
									onClick={() => logout()}
									className="ml-3 px-3 py-2 text-sm text-gray-700 hover:text-blue-600"
								>
									Sign Out
								</button>
							</div>
						) : (
							<div>
								<Link
									href="/login"
									className="px-3 py-2 text-sm text-gray-700 hover:text-blue-600"
								>
									Sign In
								</Link>
								<Link
									href="/register"
									className="ml-3 px-3 py-2 rounded-md text-sm font-medium text-white bg-blue-600 hover:bg-blue-700"
								>
									Sign Up
								</Link>
							</div>
						)}
					</div>

					<div className="-mr-2 flex items-center sm:hidden">
						<button
							onClick={() =>
								setIsMobileMenuOpen(!isMobileMenuOpen)
							}
							className="inline-flex items-center justify-center p-2 rounded-md text-gray-400 hover:text-gray-500 hover:bg-gray-100 focus:outline-none focus:ring-2 focus:ring-inset focus:ring-blue-500"
						>
							<span className="sr-only">Open main menu</span>
							{isMobileMenuOpen ? (
								<XMarkIcon
									className="block h-6 w-6"
									aria-hidden="true"
								/>
							) : (
								<Bars3Icon
									className="block h-6 w-6"
									aria-hidden="true"
								/>
							)}
						</button>
					</div>
				</div>
			</div>

			{/* Mobile menu */}
			{isMobileMenuOpen && (
				<div className="sm:hidden">
					<div className="pt-2 pb-3 space-y-1">
						{navItems.map((item) => (
							<Link
								key={item.href}
								href={item.href}
								className={`block pl-3 pr-4 py-2 border-l-4 text-base font-medium ${
									pathname === item.href
										? "border-blue-500 text-blue-700 bg-blue-50"
										: "border-transparent text-gray-600 hover:bg-gray-50 hover:border-gray-300 hover:text-gray-800"
								}`}
								onClick={() => setIsMobileMenuOpen(false)}
							>
								{item.name}
							</Link>
						))}
					</div>

					<div className="pt-4 pb-3 border-t border-gray-200">
						{isAuthenticated ? (
							<div>
								<div className="px-4 py-2 text-sm text-gray-700">
									Hi, {user?.displayName}
								</div>
								<button
									onClick={() => {
										logout();
										setIsMobileMenuOpen(false);
									}}
									className="block w-full text-left px-4 py-2 text-base font-medium text-gray-600 hover:bg-gray-100"
								>
									Sign Out
								</button>
							</div>
						) : (
							<div className="space-y-1">
								<Link
									href="/login"
									className="block px-4 py-2 text-base font-medium text-gray-600 hover:bg-gray-100"
									onClick={() => setIsMobileMenuOpen(false)}
								>
									Sign In
								</Link>
								<Link
									href="/register"
									className="block px-4 py-2 text-base font-medium text-gray-600 hover:bg-gray-100"
									onClick={() => setIsMobileMenuOpen(false)}
								>
									Sign Up
								</Link>
							</div>
						)}
					</div>
				</div>
			)}
		</nav>
	);
}
