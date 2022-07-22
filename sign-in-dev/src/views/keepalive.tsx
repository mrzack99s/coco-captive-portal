import { Chip } from 'primereact/chip'
import { Button } from 'primereact/button'
import { useApiConnector } from '../utils/api-connector';
import { useNavigate } from 'react-router';
import { useCaptivePortalProperties, useLang, usePageWaiting, useToast } from '../utils/properties';
import { useEffect, useState } from 'react';

/* eslint-disable */
const KeepAlive = () => {
    const apiInstance = useApiConnector();
    const navigate = useNavigate()
    const toast = useToast();
    const captivePortalProperties = useCaptivePortalProperties();
    const lang = useLang();
    const [htmlReady, setHtmlReady] = useState(false);
    const setPageWaiting = usePageWaiting();

    const signOut = () => {
        apiInstance.api.signOut()
            .then(() => {
                toast.current.show({ severity: 'success', summary: 'Success', detail: lang == "en" ? "Signed out" : "ออกจากระบบแล้ว", life: 3000 });
                navigate("/")
            })
            .catch(() => {
                toast.current.show({ severity: 'error', summary: 'Error', detail: lang == "en" ? "Sign out failed" : "ออกจากระบบไม่สำเร็จ", life: 3000 });
            })
    }

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
                    <div className="grid grid-nogutter p-fluid h-screen bg-bluegray-50">
                        <div className="hidden lg:flex col-7 grid-nogutter overflow-hidden p-0">
                            <img
                                src={captivePortalProperties.html!.background_file_name ? require(`../assets/${captivePortalProperties.html!.background_file_name!}`) : require("../assets/bg-01.jpg")}
                                className="relative ml-auto block w-full h-full"
                                alt='Backgroung'
                                style={{
                                    padding: 0,
                                    objectFit: 'cover'
                                }}
                            />
                        </div>
                        <div className="col grid-nogutter p-8 text-left flex align-items-center justify-content-center">
                            <section className='text-center'>
                                <img className="relative mx-auto block"
                                    alt='Logo'
                                    src={captivePortalProperties.html!.logo_file_name ? require(`../assets/${captivePortalProperties.html!.logo_file_name!}`) : require("../assets/logo.png")}
                                    style={{ maxWidth: '160px' }} />
                                <span className="block text-center text-6xl text-primary font-bold mb-1">
                                    {lang === "en" &&
                                        <>
                                            {captivePortalProperties.html!.en_title_name ? captivePortalProperties.html!.en_title_name : "Captive Portal"}
                                        </>
                                    }
                                    {lang === "th" &&
                                        <>
                                            {captivePortalProperties.html!.th_title_name ? captivePortalProperties.html!.th_title_name : "Captive Portal"}
                                        </>
                                    }
                                </span>
                                <p className="mt-0 text-center mb-4 text-700 line-height-3">
                                    {lang === "en" &&
                                        <>
                                            {captivePortalProperties.html!.en_sub_title ? captivePortalProperties.html!.en_sub_title : "Internet access authentication"}
                                        </>
                                    }
                                    {lang === "th" &&
                                        <>
                                            {captivePortalProperties.html!.th_sub_title ? captivePortalProperties.html!.th_sub_title : "Internet access authentication"}
                                        </>
                                    }
                                </p>

                                <Chip label={lang === "en" ? "Internet Access" : "เชื่อมต่ออินเทอร์เน็ต"} icon="pi pi-globe" className="mr-2 mb-2 py-4 px-6 font-bold text-3xl bg-green-800 text-white " />

                                <span className="p-buttonset mt-5">
                                    <Button onClick={() => {
                                        window.open("https://www.google.com", "_blank")
                                    }} className='p-button-sm' label={lang === "en" ? "Open Google" : "เปิด Google"} icon="pi pi-check" />
                                    <Button onClick={() => {
                                        let params = `scrollbars=no,resizable=no,status=no,location=no,toolbar=no,menubar=no,
                                    width=600,height=700,left=-1000,top=-1000`;

                                        window.open('/keepalive', 'COCO-Captive-Portal', params);
                                    }} className='p-button-sm' label={lang === "en" ? "Open keepalive popup" : "เปิดป๊อปอัพหน้าสถานะ"} icon="pi pi-external-link" />
                                </span>
                                <Button onClick={() => { signOut() }} className='p-button-sm mt-1 p-button-danger' label={lang === "en" ? "Sign-out" : "ออกจากระบบ"} icon="pi pi-times" />
                            </section>
                        </div>
                    </div>
                </>
            }

        </>
    );
}


export default KeepAlive;