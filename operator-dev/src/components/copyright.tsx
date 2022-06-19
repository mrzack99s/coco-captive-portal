export const Copyright = () => {
    const year = new Date().getFullYear()
    return (
        <div className="flex align-items-center justify-content-center text-md text-gray-400" style={{ position: 'absolute', bottom: '30px' }}>
            &copy; {year} CoCo Captive Portal.
        </div>
    )
}