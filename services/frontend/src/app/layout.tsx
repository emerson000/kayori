import { BodyWithSidebar } from '../components/navBar'
import './global.css'

export default function RootLayout({
  children,
}: {
  children: React.ReactNode
}) {
  return (
    <html lang="en">
      <body>
        <BodyWithSidebar>
          {children}
        </BodyWithSidebar>
      </body>
    </html>
  )
}