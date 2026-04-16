import type { Metadata } from "next";
import localFont from "next/font/local";
import "./globals.css";
import RefineProvider from "@/providers/refine-provider";

const jakarta = localFont({
  src: "./fonts/GeistVF.woff",
  variable: "--font-jakarta",
  weight: "300 400 500 600 700",
});

export const metadata: Metadata = {
  title: "Integrated Procurement Dashboard",
  description: "Integrated Procurement & Logistics Dashboard",
  manifest: "/manifest.webmanifest",
  icons: {
    icon: "/icon.svg",
    shortcut: "/icon.svg",
    apple: "/icon.svg",
  },
};

export const dynamic = "force-dynamic";

export default function RootLayout({
  children,
}: Readonly<{
  children: React.ReactNode;
}>) {
  return (
    <html lang="en">
      <body
        className={`${jakarta.variable} antialiased`}
      >
        <RefineProvider>{children}</RefineProvider>
      </body>
    </html>
  );
}
