/* eslint-disable */
import React, { useEffect, useState } from "react";
import { InputText } from "primereact/inputtext";
import { Button } from "primereact/button";
import { useApiConnector } from "../utils/api-connector";
import { useNavigate } from "react-router";
import { Copyright } from "../components/copyright"
import {
    useHtmlProperties,
    useInitSecret,
    useLang,
    usePageWaiting,
    useRedirectURL,
    useToast,
} from "../utils/properties";
import { Password } from "primereact/password";
import { useLocation } from "react-router-dom";

const SignIn = () => {
    const apiInstance = useApiConnector();
    const [username, setUsername] = useState("");
    const [password, setPassword] = useState("");
    const navigate = useNavigate();
    const location = useLocation();
    const toast = useToast();
    const htmlProperties = useHtmlProperties();
    const [initSecret, setInitSecret] = useInitSecret();
    const lang = useLang();
    const setPageWaiting = usePageWaiting();
    const [redirectUrl, setRedirectUrl] = useRedirectURL();

    const goLogin = () => {
        apiInstance.api
            .authentication({
                username: username,
                password: password,
                secret: initSecret,
            })
            .then((res) => res.data)
            .then((res) => {
                toast.current.show({
                    severity: "success",
                    summary: "Success",
                    detail: lang == "en" ? "Signed" : "เข้าสู่ระบบแล้ว",
                    life: 3000,
                });
                setRedirectUrl(res.redirect_url!);
                if (res.redirect_url && res.redirect_url !== "") {
                    let params = `scrollbars=no,resizable=no,status=no,location=no,toolbar=no,menubar=no,
                                    width=600,height=700,left=-1000,top=-1000`;

                    window.open("/keepalive", "COCO-Captive-Portal", params);
                    window.location.href = res.redirect_url;
                } else {
                    navigate("/keepalive");
                }
                setPageWaiting(false);
            })
            .catch(() => {
                toast.current.show({
                    severity: "error",
                    summary: "Error",
                    detail: lang == "en" ? "Wrong credential!" : "ข้อมูลประจำตัวไม่ถูกต้อง",
                    life: 3000,
                });
                setPageWaiting(false);
            });
    }
    const verify = (event: { preventDefault: () => void }) => {
        event.preventDefault();
        setPageWaiting(true);
        apiInstance.api
            .isExistInitializeSecret({
                secret: initSecret,
            })
            .then(() => {
                goLogin()
            })
            .catch(() => {
                apiInstance.api
                    .initialize()
                    .then(res => res.data)
                    .then((res) => {
                        setInitSecret(res.secret!)
                        goLogin()
                    })
                    .catch(() => {
                        toast.current.show({
                            severity: "error",
                            summary: "Error",
                            detail: lang == "en" ? "Not found initial secret, reload in 3 seconds" : "ไม่พบเซสชันเริ่มต้น, รีโหลดตัวเองใน 3 วินาที",
                            life: 3000,
                        });
                        setTimeout(() => {
                            window.location.reload()
                        }, 3000)

                    })
            })

    };

    return (
        <>
            <div>
                <div className="grid grid-nogutter p-fluid h-screen bg-bluegray-50">
                    <div className="hidden lg:flex col-7 grid-nogutter overflow-hidden p-0">
                        <img
                            src={
                                htmlProperties.background_file_name
                                    ? require(`../assets/${htmlProperties.background_file_name!}`)
                                    : require("../assets/bg-01.jpg")
                            }
                            className="relative ml-auto block w-full h-full"
                            style={{
                                padding: 0,
                                objectFit: "cover",
                            }}
                        />
                    </div>
                    <div className="col grid-nogutter p-8 text-left flex align-items-center justify-content-center">
                        <section>
                            <img
                                className="relative mx-auto block"
                                src={
                                    htmlProperties.logo_file_name
                                        ? require(`../assets/${htmlProperties.logo_file_name!}`)
                                        : require("../assets/logo.png")
                                }
                                style={{ maxWidth: "160px" }}
                            />
                            <span className="block text-center text-6xl text-primary font-bold mb-1">
                                {lang === "en" && (
                                    <>
                                        {htmlProperties.en_title_name
                                            ? htmlProperties.en_title_name
                                            : "Captive Portal"}
                                    </>
                                )}
                                {lang === "th" && (
                                    <>
                                        {htmlProperties.th_title_name
                                            ? htmlProperties.th_title_name
                                            : "Captive Portal"}
                                    </>
                                )}
                            </span>
                            <p className="mt-0 text-center mb-4 text-700 line-height-3">
                                {lang === "en" && (
                                    <>
                                        {htmlProperties.en_sub_title
                                            ? htmlProperties.en_sub_title
                                            : "Internet access authentication"}
                                    </>
                                )}
                                {lang === "th" && (
                                    <>
                                        {htmlProperties.th_sub_title
                                            ? htmlProperties.th_sub_title
                                            : "Internet access authentication"}
                                    </>
                                )}
                            </p>
                            <form onSubmit={verify}>
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
};

export default SignIn;
