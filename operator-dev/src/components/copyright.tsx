export const Copyright = () => {
    const year = new Date().getFullYear()
    return (
        <div className="flex align-items-center justify-content-center text-md text-gray-400 mt-8">
            &copy; {year} CoCo Captive Portal.
        </div>
    )
}