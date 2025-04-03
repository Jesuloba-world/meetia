import type { Metadata } from "next";
import { Geist_Mono, Inter } from "next/font/google";
import "./globals.css";
import { Providers } from "./providers";

const InterFont = Inter({
	variable: "--font-inter-sans",
	subsets: ["latin"],
});

const geistMono = Geist_Mono({
	variable: "--font-geist-mono",
	subsets: ["latin"],
});

export const metadata: Metadata = {
	title: "Meetia - Video Conferencing",
	description: "A Google Meet clone built with Next.js and Go",
};

export default function RootLayout({
	children,
}: Readonly<{
	children: React.ReactNode;
}>) {
	return (
		<html lang="en">
			<body
				className={`${InterFont.variable} ${geistMono.variable} antialiased font-sans`}
			>
				<Providers>{children}</Providers>
			</body>
		</html>
	);
}
