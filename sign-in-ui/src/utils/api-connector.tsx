import { Api, ApiConfig } from "../api"
import { useContext, createContext } from "react";
import { useCookies } from "react-cookie";

type ContextProps = {
    instance: Api<unknown>
}

export const ApiContext = createContext<ContextProps | null>(null);

interface ApiPros {
    children: any
}


const ApiProvider: React.FC<ApiPros> = ({ children }) => {
    const apiConnector = new Api({
        baseUrl: process.env.REACT_APP_ENV === "production" ? "/api" : `${process.env.REACT_APP_API_URL}/api`,

    })
    return <ApiContext.Provider value={{ instance: apiConnector }}>{children}</ApiContext.Provider>;
}

export default ApiProvider;
export const useApiConnector = () => {
    const apiContext = useContext(ApiContext) as ContextProps
    return apiContext.instance;
}

export const useAdminApiConnector = () => {
    const [cookies, setCookie, removeCookie] = useCookies(['api-token']);
    const apiContext = useContext(ApiContext) as ContextProps
    const secureApi = new Api({
        baseUrl: apiContext.instance.baseUrl,
        baseApiParams: {
            headers: {
                "api-token": cookies["api-token"]
            }
        }
    })
    return secureApi;
}
