import { useContext, createContext, useState, useEffect, useRef } from "react";
import { useNavigate, useLocation } from "react-router-dom";
import { useAdminApiConnector } from "./api-connector";
import { Toast } from "primereact/toast";
import { TypesHTMLType } from "../api";
import { Button } from 'primereact/button';
import { useCookies } from "react-cookie";
import { ProgressSpinner } from 'primereact/progressspinner';
import { ConfirmDialog } from "primereact/confirmdialog";

type ContextProps = {
    issue: string
    initSecret: string
    lang: string
    redirectUrl: string
    setRedirectUrl: React.Dispatch<React.SetStateAction<string>>
    setInitSecret: React.Dispatch<React.SetStateAction<string>>
    setWaiting: React.Dispatch<React.SetStateAction<boolean>>
    htmlProperties: TypesHTMLType
    toastRef: React.MutableRefObject<any>
}

export const AppPropertiesContext = createContext<ContextProps | null>(null);

interface AppProperties {
    children: any
}

const AppPropertiesProvider: React.FC<AppProperties> = ({ children }) => {
    // eslint-disable-next-line
    const [issue, setIssue] = useState("")
    const [redirectUrl, setRedirectUrl] = useState("")
    const [initSecret, setInitSecret] = useState("")
    // eslint-disable-next-line
    const [htmlProperties, setHtmlProperties] = useState({} as TypesHTMLType)
    const navigate = useNavigate()
    const location = useLocation();
    const [pageWaiting, setPageWaiting] = useState(false)
    const apiInstance = useAdminApiConnector();
    const toast = useRef({} as any)
    const [, , removeCookies] = useCookies(["api-token"])

    const checkCredential = (path: string) => {
        apiInstance.api.admSigned()
            .then(() => {
                if (path === "/sign-in") {
                    navigate("/overview")
                }
            })
            .catch((err) => {
                if ((err.error.msg as string).trim() === "your network is not authorized") {
                    removeCookies("api-token")
                    navigate("/unauthorized")
                }
                else {
                    removeCookies("api-token")
                    navigate("/sign-in")
                }
            })
    }


    useEffect(() => {

        if (!location.pathname.includes("/unauthorized")) {
            checkCredential(location.pathname)
        }
        // eslint-disable-next-line react-hooks/exhaustive-deps
    }, [location.pathname])


    const [lang, setLang] = useState('en');

    return (

        <AppPropertiesContext.Provider value={{ redirectUrl: redirectUrl, setRedirectUrl: setRedirectUrl, lang: lang, setWaiting: setPageWaiting, issue: issue, initSecret: initSecret, setInitSecret: setInitSecret, htmlProperties: htmlProperties, toastRef: toast }}>
            <Toast ref={toast} />
            <ConfirmDialog />
            {pageWaiting &&
                <div className="bg-gray-50 h-screen w-screen" style={{ position: 'fixed', zIndex: 1, opacity: 0.8 }}>
                    <div className="flex align-items-center h-screen justify-content-center">
                        <span className="text-center">
                            <p>
                                <ProgressSpinner style={{ width: '100px', height: '100px' }} strokeWidth="5" fill="var(--surface-ground)" animationDuration=".5s" />
                            </p>
                            <p className="text-2xl">{lang === "en" ? "Waiting" : "กรุณารอสักครู่"}</p>
                        </span>
                    </div>


                </div>
            }

            {location.pathname.includes("/sign-in") &&
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
            }
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

export const useHtmlProperties = () => {
    const appContext = useContext(AppPropertiesContext) as ContextProps
    return appContext.htmlProperties;
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