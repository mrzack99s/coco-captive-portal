import React, { useState } from "react";
import { Menubar } from "primereact/menubar";
import { Button } from "primereact/button";
import { Sidebar } from "primereact/sidebar";
import { Avatar } from "primereact/avatar";
import { Divider } from "primereact/divider";
import './custom.css'
import { useCookies } from 'react-cookie';
import { useNavigate } from 'react-router-dom';
import { useAdminApiConnector } from "../utils/api-connector";
import { useToast } from "../utils/properties";

const items = [
    {
        label: "Management",
        icon: "pi pi-building",
        items: [
            {
                label: "Session Monitor",
                icon: "mdi mdi-monitor-eye",
                url: "/session-monitor"
            },
            {
                separator: true,
            },
            {
                label: "Policy and Objects",
                icon: "mdi mdi-cube-outline",
                url: "/policy-and-objects"
            },
            {
                label: "Config",
                icon: "pi pi-sliders-h",
                url: "/config"
            },
        ],
    },
];

export default () => {
    const [cookies, setCookie, removeCookie] = useCookies(['role', 'api-token']);
    const navigate = useNavigate()
    const apiInstance = useAdminApiConnector()
    const toast = useToast();

    const revokeAdmin = () => {
        apiInstance.api.revokeAdministrator()
        removeCookie("api-token")
        toast.current.show({ severity: 'success', summary: 'Success', detail: `Sigend out`, life: 3000 });
        navigate("/operator/sign-in")
    }

    const end = () => (
        <>
            <Button
                icon="pi pi-user"
                className="my-1 p-0 p-button-rounded surface-ground text-color-secondary"
                aria-label="Filter"
                style={{
                    borderRadius: "100%",
                    height: "40px",
                    width: "40px"
                }}
                onClick={() => setVisibleRight(true)}
            />
        </>
    );
    const start = () => (
        <div><Button label="COCO" onClick={() => navigate("/overview")} className="p-button-text text-white text-2xl p-0 bg-blue-500 mr-4" /></div>
    );
    const [visibleRight, setVisibleRight] = useState(false);


    return (
        <>
            <Sidebar
                visible={visibleRight}
                position="right"
                onHide={() => setVisibleRight(false)}
            >
                <div className="text-center">
                    <Avatar icon="pi pi-user" className="mr-2" size="xlarge" shape="circle" />
                    <p>Administrator</p>
                </div>
                <Divider align="left" type="dashed" className="mt-4 mb-0">
                    <span className="text-sm">Preferences</span>
                </Divider>
                <div className="mt-4 text-center">
                    <Button label="Sign Out" className="text-xs p-button-danger" icon="pi pi-sign-out" onClick={() => {
                        revokeAdmin()
                    }} />
                </div>
            </Sidebar>

            < Menubar
                className="py-0 px-3 bg-blue-500"
                start={start}
                model={items}
                end={end}
                style={{ width: "100%", position: "fixed", zIndex: 100 }}
            />

        </>
    );
};