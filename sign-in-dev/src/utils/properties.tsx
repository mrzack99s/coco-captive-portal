import { useContext, createContext, useState, useEffect, useRef } from "react";
import { useNavigate, useLocation } from "react-router-dom";
import { useApiConnector } from "./api-connector";
import { Toast } from "primereact/toast";
import { TypesCaptivePortalConfigFundamentalType } from "../api";
import { Button } from 'primereact/button';
import { ProgressSpinner } from 'primereact/progressspinner';

type ContextProps = {
    issue: string
    initSecret: string
    lang: string
    redirectUrl: string
    setRedirectUrl: React.Dispatch<React.SetStateAction<string>>
    setInitSecret: React.Dispatch<React.SetStateAction<string>>
    setWaiting: React.Dispatch<React.SetStateAction<boolean>>
    captivePortalProperties: TypesCaptivePortalConfigFundamentalType
    toastRef: React.MutableRefObject<any>
}

export const AppPropertiesContext = createContext<ContextProps | null>(null);

interface AppProperties {
    children: any
}

const AppPropertiesProvider: React.FC<AppProperties> = ({ children }) => {

    const [issue, setIssue] = useState("")
    const [redirectUrl, setRedirectUrl] = useState("")
    const [defaultLangSet, setDefaultLangSet] = useState(false)
    const [initSecret, setInitSecret] = useState("")
    const [captivePortalProperties, setCaptivePortalProperties] = useState({} as TypesCaptivePortalConfigFundamentalType)
    const navigate = useNavigate()
    const location = useLocation();
    const [pageWaiting, setPageWaiting] = useState(false)
    const apiInstance = useApiConnector();
    const toast = useRef({} as any)

    const checkCredential = () => {
        apiInstance.api.signed()
            .then(res => res.data)
            .then(res => {
                setIssue(res.issue!)
                navigate("/keepalive")
            })
            .catch(() => {
                navigate("/sign-in")
            })
    }

    const getHtmlProps = () => {
        apiInstance.api.getCaptivePortalConfigFundamental()
            .then(res => res.data)
            .then(res => {
                setCaptivePortalProperties(res!)
                if (!defaultLangSet) {
                    setLang(res.html!.default_language!)
                    setDefaultLangSet(true)
                }

            })
    }


    useEffect(() => {
        getHtmlProps()
        checkCredential()
        if (location.pathname.includes("/sign-in")) {
            apiInstance.api.initialize()
                .then(res => res.data)
                .then(res => {
                    setInitSecret(res.secret!)
                    if (!defaultLangSet) {
                        navigate("/sign-in")
                    }
                })
        }
        // eslint-disable-next-line react-hooks/exhaustive-deps
    }, [location.pathname])


    const [lang, setLang] = useState('en');

    return (

        <AppPropertiesContext.Provider value={{ redirectUrl: redirectUrl, setRedirectUrl: setRedirectUrl, lang: lang, setWaiting: setPageWaiting, issue: issue, initSecret: initSecret, setInitSecret: setInitSecret, captivePortalProperties: captivePortalProperties, toastRef: toast }}>
            <Toast ref={toast} />
            {pageWaiting &&
                <div className="bg-gray-50 h-screen w-screen" style={{ position: 'fixed', zIndex: 10000, opacity: 0.8 }}>
                    <div className="flex align-items-center h-screen justify-content-center">
                        <span className="text-center">
                            <div>
                                <ProgressSpinner style={{ width: '100px', height: '100px' }} strokeWidth="5" fill="var(--surface-ground)" animationDuration=".5s" />
                            </div>
                            <p className="text-2xl">{lang === "en" ? "Waiting" : "กรุณารอสักครู่"}</p>
                        </span>
                    </div>


                </div>
            }



            <span className="p-buttonset" style={{
                position: "absolute",
                top: "30px",
                right: "3%"
            }}>
                <Button label="English"
                    onClick={() => {
                        setLang('en')
                    }}
                    className={`p-button-xs pt-1 px-3 text-sm ${lang !== "en" ? "p-button-outlined" : ""}`} />
                <Button label="ภาษาไทย"
                    onClick={() => {
                        setLang('th')
                    }}
                    className={`p-button-xs pt-1 px-3 text-sm ${lang !== "th" ? "p-button-outlined" : ""}`} />
            </span>


            {children}
        </AppPropertiesContext.Provider>
    );
}

export default AppPropertiesProvider;
export const useIssue = () => {
    const appContext = useContext(AppPropertiesContext) as ContextProps
    return [appContext.issue,];
}

export const useToast = () => {
    const appContext = useContext(AppPropertiesContext) as ContextProps
    return appContext.toastRef;
}

export const useCaptivePortalProperties = () => {
    const appContext = useContext(AppPropertiesContext) as ContextProps
    return appContext.captivePortalProperties;
}

export const useInitSecret = (): [string, React.Dispatch<React.SetStateAction<string>>] => {
    const appContext = useContext(AppPropertiesContext) as ContextProps
    return [appContext.initSecret, appContext.setInitSecret];
}

export const useRedirectURL = (): [string, React.Dispatch<React.SetStateAction<string>>] => {
    const appContext = useContext(AppPropertiesContext) as ContextProps
    return [appContext.redirectUrl, appContext.setRedirectUrl];
}

export const useLang = (): string => {
    const appContext = useContext(AppPropertiesContext) as ContextProps
    return appContext.lang;
}

export const usePageWaiting = (): React.Dispatch<React.SetStateAction<boolean>> => {
    const appContext = useContext(AppPropertiesContext) as ContextProps
    return appContext.setWaiting;
}