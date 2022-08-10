import { useCaptivePortalProperties, useLang, usePageWaiting } from '../utils/properties';
import { useEffect, useState } from 'react';

/* eslint-disable */
const UnAuthorized = () => {
    const captivePortalProperties = useCaptivePortalProperties();
    const lang = useLang();
    const setPageWaiting = usePageWaiting();
    const [htmlReady, setHtmlReady] = useState(false);

    useEffect(() => {
        if (captivePortalProperties.html) {
            setPageWaiting(false)
            setHtmlReady(true)
        } else {
            setPageWaiting(true)
            setHtmlReady(false)
        }
    }, [captivePortalProperties.html])

    return (
        <>
            {htmlReady &&
                <>
                    <div className="p-fluid h-screen bg-bluegray-50">

                        <div className="h-screen text-center flex align-items-center justify-content-center">
                            <section>
                                <img
                                    className="relative mx-auto block"
                                    src={
                                        captivePortalProperties.html!.logo_file_name
                                            ? require(`../assets/${captivePortalProperties.html!.logo_file_name!}`)
                                            : require("../assets/logo.png")
                                    }
                                    style={{ maxWidth: "160px" }}
                                />

                                <span className="block text-6xl text-primary font-bold mb-1">
                                    {lang === "en" && (
                                        <>
                                            {captivePortalProperties.html!.en_title_name
                                                ? captivePortalProperties.html!.en_title_name
                                                : "CoCo Captive Portal"}
                                        </>
                                    )}
                                    {lang === "th" && (
                                        <>
                                            {captivePortalProperties.html!.th_title_name
                                                ? captivePortalProperties.html!.th_title_name
                                                : "CoCo Captive Portal"}
                                        </>
                                    )}
                                </span>
                                <p className="mt-3 block border-round-3xl mb-4 text-3xl text-bold line-height-3 bg-gray-50 p-5 text-red-500">
                                    Your network is not authorized !!
                                </p>
                            </section>
                        </div>
                    </div>
                </>
            }
        </>
    );
}


export default UnAuthorized;