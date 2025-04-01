import Link from "next/link";

const navItems = [
  { label: "Home", href: "/" },
  { label: "Projects", href: "/projects" },
]

export function NavBar() {
  return <div className="navbar bg-base-100">
    <div className="navbar-start">
      <label htmlFor="my-drawer-3" aria-label="open sidebar" className="btn btn-square btn-ghost lg:hidden">
        <svg
          xmlns="http://www.w3.org/2000/svg"
          fill="none"
          viewBox="0 0 24 24"
          className="inline-block h-6 w-6 stroke-current">
          <path
            strokeLinecap="round"
            strokeLinejoin="round"
            strokeWidth="2"
            d="M4 6h16M4 12h16M4 18h16"></path>
        </svg>
      </label>
      <Link href="/" className="btn btn-ghost text-xl">Kayori</Link>
      <ul className="menu menu-horizontal px-1 hidden lg:flex">
        {navItems && navItems.map((item) => <li key={item.href}>
          <Link href={item.href}>{item.label}</Link>
        </li>)}
      </ul>
    </div>
  </div >
}

export function BodyWithSidebar(props: { children: React.ReactNode }) {
  return <div className="drawer">
    <input id="my-drawer-3" type="checkbox" className="drawer-toggle" />
    <div className="drawer-content flex flex-col">
      <NavBar />
      {props.children}
    </div>
    <div className="drawer-side">
      <label htmlFor="my-drawer-3" aria-label="close sidebar" className="drawer-overlay"></label>
      <ul className="menu bg-base-200 min-h-full w-80 p-4">
        {navItems && navItems.map((item) => <li key={item.href}>
          <Link href={item.href}>{item.label}</Link>
        </li>)}
      </ul>
    </div>
  </div>
}