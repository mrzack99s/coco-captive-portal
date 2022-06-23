import { useState } from 'react';
import { InputText } from 'primereact/inputtext'
import { Button } from 'primereact/button'
import { useApiConnector } from '../utils/api-connector';
import { useNavigate } from 'react-router';
import { useHtmlProperties, useLang, usePageWaiting, useToast } from '../utils/properties';
import { Password } from 'primereact/password';
import { useCookies } from 'react-cookie';
import { Copyright } from '../components/copyright';

/* eslint-disable */
function SignIn() {
    const apiInstance = useApiConnector();
    const [username, setUsername] = useState("");
    const [password, setPassword] = useState("");
    const navigate = useNavigate()
    const toast = useToast();
    const htmlProperties = useHtmlProperties();
    const [cookies, setCookies] = useCookies(["api-token"])
    const lang = useLang();
    const setPageWaiting = usePageWaiting()

    const goLogin = (event: { preventDefault: () => void; }) => {
        event.preventDefault()
        setPageWaiting(true)
        apiInstance.api.checkIsAdministrator({
            username: username,
            password: password,
        })
            .then(res => res.data)
            .then(res => {
                console.log(res)
                let now = new Date();
                now.setTime(now.getTime() + 1 * 60 * 60 * 1000);
                setCookies("api-token", res.api_token, {
                    expires: now,
                })
                toast.current.show({ severity: 'success', summary: 'Success', detail: "Signed", life: 3000 });
                setPageWaiting(false)
                navigate("/operator/monitor")
            })
            .catch(() => {
                toast.current.show({ severity: 'error', summary: 'Error', detail: "Wrong credential!", life: 3000 });
                setPageWaiting(false)
            })
    }

    return (
        <>

            <div>
                <div className="grid grid-nogutter p-fluid h-screen bg-bluegray-50">
                    <div className="hidden lg:flex col-7 grid-nogutter overflow-hidden p-0">
                        <img
                            src={htmlProperties.background_file_name ? require(`../assets/${htmlProperties.background_file_name!}`) : require("../assets/bg-01.jpg")}
                            className="relative ml-auto block w-full h-full"
                            style={{
                                padding: 0,
                                objectFit: 'cover'
                            }}
                        />
                    </div>
                    <div className="col grid-nogutter p-8 text-left flex align-items-center justify-content-center">
                        <section>
                            <img className="relative mx-auto block"
                                src={htmlProperties.logo_file_name ? require(`../assets/${htmlProperties.logo_file_name!}`) : require("../assets/logo.png")}
                                style={{ maxWidth: '160px' }} />
                            <span className="block text-center text-6xl text-primary font-bold mb-1">
                                {lang === "en" ? "CoCo Operator" : "โคโค่โอเปอเรเตอร์"}
                            </span>
                            <p className="mt-0 text-center mb-4 text-700 line-height-3">
                                {lang === "en" ? "Configuration manager and session monitor" : "ตัวจัดการคอนฟิกและเซสชันมอนิเตอร์"}
                            </p>
                            <form onSubmit={goLogin}>
                                <div className="p-inputgroup  m-2">
                                    <span className="p-inputgroup-addon">
                                        <i className="pi pi-user" />
                                    </span>
                                    <InputText
                                        value={username}
                                        onChange={(e) => setUsername(e.target.value)}
                                        placeholder={lang === "en" ? "Username" : "ชื่อผู้ใช้"}
                                        required
                                    />
                                </div>

                                <div className="p-inputgroup m-2">
                                    <span className="p-inputgroup-addon">
                                        <i className="pi pi-key" />
                                    </span>
                                    <Password
                                        value={password}
                                        onChange={(e) => setPassword(e.target.value)}
                                        placeholder={lang === "en" ? "Password" : "รหัสผ่าน"}
                                        feedback={false}
                                        required
                                        toggleMask
                                    />
                                </div>

                                <Button
                                    label={lang === "en" ? "Sign In" : "ลงชื่อเข้าใช้"}
                                    type="submit"
                                    className="p-button m-2"
                                />
                            </form>
                            <Copyright />
                        </section>
                    </div>
                </div>
            </div>

        </>
    );
}

export default SignIn;

