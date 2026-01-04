import type { Metadata } from "next";
import "./globals.css";

export const metadata: Metadata = {
  title: "ShadowNet - P2P Mesh VPN Dashboard",
  description: "Monitor and manage your ShadowNet peer-to-peer mesh VPN network",
};

export default function RootLayout({
  children,
}: Readonly<{
  children: React.ReactNode;
}>) {
  return (
    <html lang="en" className="dark">
      <body className="antialiased bg-black text-white">
        {children}
      </body>
    </html>
  );
}
