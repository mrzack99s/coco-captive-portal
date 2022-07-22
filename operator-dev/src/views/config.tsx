import Navbar from "../components/navbar";
import Config from "../components/configs";
import { Copyright } from "../components/copyright"

const ConfigView = () => {

    return (
        <div className="mb-5">
            <Navbar />
            <div className="grid grid-nogutter m-0" style={{ position: "relative", top: "65px" }}>
                <div className="col hidden lg:inline grid-nogutter"></div>
                <div className="col-12 lg:col-8 grid-nogutter" style={{}}>
                    <Config />
                </div>
                <div className="col hidden lg:inline grid-nogutter"></div>
            </div>
            <Copyright />
        </div>
    );
};

export default ConfigView;