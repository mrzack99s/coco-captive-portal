import { useLang, useHtmlProperties } from '../utils/properties';

/* eslint-disable */
const UnAuthorized = () => {
    const lang = useLang();
    const htmlProperties = useHtmlProperties();

    return (
        <>

            <div className="p-fluid h-screen bg-bluegray-50">

                <div className="h-screen text-center flex align-items-center justify-content-center">
                    <section>
                        <img className="relative mx-auto block"
                            src={htmlProperties.logo_file_name ? require(`../assets/${htmlProperties.logo_file_name!}`) : require("../assets/logo.png")}
                            style={{ maxWidth: '160px' }} />
                        <span className="block text-center text-6xl text-primary font-bold mb-1">
                            {lang === "en" ? "CoCo Operator" : "โคโค่โอเปอเรเตอร์"}
                        </span>
                        <p className="mt-3 block border-round-3xl mb-4 text-3xl text-bold line-height-3 bg-gray-50 p-5 text-red-500">
                            Your network is not authorized !!
                        </p>
                    </section>
                </div>
            </div>
        </>
    );
}


export default UnAuthorized;