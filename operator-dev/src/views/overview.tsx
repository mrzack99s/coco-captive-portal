import Navbar from "../components/navbar";
import { useCookies } from 'react-cookie';
import Monitor from '../components/monitor'
import { useEffect } from "react";
import { useAdminApiConnector } from "../utils/api-connector";
import Overview from "../components/overview";
export default () => {
    const [cookies, setCookie, removeCookie] = useCookies(['role']);

    return (
        <div>
            <Navbar />
            <div className="grid grid-nogutter m-0" style={{ position: "relative", top: "65px" }}>
                <div className="col hidden lg:inline grid-nogutter"></div>
                <div className="col-12 lg:col-8 grid-nogutter" style={{}}>
                    <Overview />
                </div>
                <div className="col hidden lg:inline grid-nogutter"></div>
            </div>
        </div>
    );
};