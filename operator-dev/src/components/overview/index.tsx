import { useEffect, useState } from 'react';
import { useAdminApiConnector } from '../../utils/api-connector';
import { useNavigate } from 'react-router';
import { useToast } from '../../utils/properties';
import { TypesConfigType, TypesExtendConfigType, TypesSessionType } from '../../api';
import { Copyright } from '../copyright';
import VirtualEditor from '../veditor';
import { confirmDialog, ConfirmDialog } from 'primereact/confirmdialog';
import { Chip } from 'primereact/chip';
import { Divider } from 'primereact/divider';
import YAML from 'yaml'
import { useClipboard } from '../../utils/clipboard';
import { Button } from 'primereact/button';

const Overview = () => {
    /* eslint-disable */
    const [countSession, setCountSession] = useState("")
    const [netIntfUsage, setNetInfUsage] = useState({
        "secure_interface": {
            "rx": "",
            "tx": "",
        },
        "egress_interface": {
            "rx": "",
            "tx": "",
        },
    })
    const [initTrig, setInitTrig] = useState(false)
    const [config, setConfig] = useState({} as TypesExtendConfigType)
    const [realtimeLoading, setRealtimeLoading] = useState(true)
    const navigate = useNavigate()
    const apiInstance = useAdminApiConnector()
    const toast = useToast();
    const clipboard = useClipboard();

    useEffect(() => {
        if (!initTrig) {
            setInterval(() => {
                setRealtimeLoading(true)
                apiInstance.api.netInterfacesBytesUsage()
                    .then(res => res.data)
                    .then(res => {
                        setNetInfUsage({
                            "secure_interface": {
                                "rx": res.secure_interface.rx,
                                "tx": res.secure_interface.tx,
                            },
                            "egress_interface": {
                                "rx": res.egress_interface.rx,
                                "tx": res.egress_interface.tx,
                            },
                        })
                        setRealtimeLoading(false)
                    })
            }, 3000)
        }
        setInitTrig(true)
    }, [initTrig])

    useEffect(() => {
        apiInstance.api.countAllSession()
            .then(res => res.data)
            .then(res => {
                setCountSession(res)
            })

        apiInstance.api.getConfig()
            .then(res => res.data)
            .then(res => {
                setConfig(res)
            })
    }, [])

    const copy = (msg: string) => {
        clipboard(msg)
            .then(() => {
                toast.current.show({ severity: 'success', summary: 'Success', detail: `Copy ${msg} to clipboard`, life: 3000 });
            })
            .catch(() => {
                toast.current.show({ severity: 'error', summary: 'Error', detail: `Copy ${msg} to clipboard failed`, life: 3000 });
            })
    }


    return (
        <>
            <div className="grid px-2">
                <div className="col-12 md:col-4 lg:col-4">
                    <div className="surface-0 p-3 border-1 border-50 border-round" style={{ height: '130px' }}>
                        <div className="flex justify-content-between mb-3">
                            <div>
                                <span className="block text-500 font-medium mb-3">Sessions</span>
                                <div className="text-900 font-medium text-xl">

                                    <Chip
                                        label={countSession ? countSession : "0"}
                                        onClick={() => {
                                            navigate("/session-monitor")
                                        }}
                                        className="hover:bg-gray-400 hover:text-white transition-duration-500 cursor-pointer" />
                                </div>
                            </div>
                            <div className="flex align-items-center justify-content-center bg-blue-100 border-round" style={{ width: '2.5rem', height: '2.5rem' }}>
                                <i className="mdi mdi-cogs text-blue-500 text-xl"></i>
                            </div>
                        </div>
                    </div>
                </div>
                <div className="col-12 md:col-4 lg:col-4">
                    <div className="surface-0 p-3 border-1 border-50 border-round" style={{ height: '130px' }}>
                        <div className="flex justify-content-between mb-3">
                            <div>
                                <span className="block text-500 font-medium mb-3">Egress Interface</span>
                                <div className="text-900 font-medium text-xl" >
                                    <Chip label={`RX: ${netIntfUsage.egress_interface.rx} MB/s`} />
                                    <Chip className='mt-1 lg:mt-0 lg:ml-2' label={`TX: ${netIntfUsage.egress_interface.tx} MB/s`} />
                                </div>
                            </div>
                            <div className="flex align-items-center justify-content-center bg-blue-100 border-round" style={{ width: '2.5rem', height: '2.5rem' }}>
                                <i className="mdi mdi-expansion-card-variant text-blue-500 text-xl"></i>
                            </div>
                        </div>
                    </div>
                </div>
                <div className="col-12 md:col-4 lg:col-4">
                    <div className="surface-0 p-3 border-1 border-50 border-round" style={{ height: '130px' }}>
                        <div className="flex justify-content-between mb-3">
                            <div>
                                <span className="block text-500 font-medium mb-3">Secure Interface</span>
                                <div className="text-900 font-medium text-xl">
                                    <Chip label={`RX: ${netIntfUsage.secure_interface.rx} MB/s`} />
                                    <Chip className='mt-1 lg:mt-0 lg:ml-2' label={`TX: ${netIntfUsage.secure_interface.tx} MB/s`} />
                                </div>
                            </div>
                            <div className="flex align-items-center justify-content-center bg-blue-100 border-round" style={{ width: '2.5rem', height: '2.5rem' }}>
                                <i className="mdi mdi-expansion-card-variant text-blue-500 text-xl"></i>
                            </div>
                        </div>
                    </div>
                </div>
            </div>
            <div className="grid px-2">
                <div className="col-12 md:col-6 lg:col-6">
                    <div className="surface-0 p-3 border-1 border-50 border-round" >
                        <div className="flex justify-content-between m1-3">
                            <div>
                                <span className="block text-500 font-medium mb-3">Informations</span>
                            </div>
                        </div>
                        <Divider />
                        <div className='flex flex-wrap align-items-center justify-content-center card-container'>
                            <div className="w-full border-round font-bold p-3 ">
                                <ul className="p-0 m-0">
                                    <li className="flex align-items-center px-2 border-300 flex-wrap">
                                        <div className={`text-500 text-right font-medium`}>
                                            Portal Mode:
                                        </div>
                                        <div className="font-normal pl-3">
                                            {config.external_portal_url ? "External" : "Built-in"}
                                        </div>
                                    </li>
                                    <li className="flex align-items-center px-2 border-300 flex-wrap">
                                        <div className={`text-500 text-right font-medium`}>
                                            Authen Mode:
                                        </div>
                                        <div className="font-normal pl-3">
                                            {config.ldap &&
                                                <>
                                                    LDAP
                                                </>
                                            }

                                            {(!config.ldap && config.radius) &&
                                                <>
                                                    Radius
                                                </>
                                            }

                                            {(!config.ldap && !config.radius) &&
                                                <>
                                                    None
                                                </>
                                            }
                                        </div>
                                    </li>
                                    <li className="flex align-items-center px-2 border-300 flex-wrap">
                                        <div className={`text-500 text-right font-medium`}>
                                            Egress Interface:
                                        </div>
                                        <div className="font-normal pl-3">
                                            {config.egress_interface}, IPv4: {config.status?.egress_ip_address}
                                            <span
                                                onClick={() => {
                                                    copy(config.status?.egress_ip_address!)
                                                }}
                                                className='ml-2 text-xs text-blue-500 hover:bg-blue-500 hover:text-white transition-duration-500 cursor-pointer border-1 px-1'>
                                                <i className="pi pi-paperclip text-xs"></i>
                                            </span>
                                        </div>
                                    </li>
                                    <li className="flex align-items-center px-2 border-300 flex-wrap">
                                        <div className={`text-500 text-right font-medium`}>
                                            Secure Interface:
                                        </div>
                                        <div className="font-normal pl-3">
                                            {config.secure_interface}, IPv4: {config.status?.secure_ip_address}
                                            <span
                                                onClick={() => {
                                                    copy(config.status?.secure_ip_address!)
                                                }}
                                                className='ml-2 text-xs text-blue-500 hover:bg-blue-500 hover:text-white transition-duration-500 cursor-pointer border-1 px-1'>
                                                <i className="pi pi-paperclip text-xs"></i>
                                            </span>
                                        </div>
                                    </li>
                                    <li className="flex align-items-center px-2 border-300 flex-wrap">
                                        <div className={`text-500 text-right font-medium`}>
                                            Session Idle Timeout:
                                        </div>
                                        <div className="font-normal pl-3">
                                            {config.session_idle} {config.session_idle! > 1 ? "minutes" : "minute"}
                                        </div>
                                    </li>
                                    <li className="flex align-items-center px-2 border-300 flex-wrap">
                                        <div className={`text-500 text-right font-medium`}>
                                            Max Concurrent:
                                        </div>
                                        <div className="font-normal pl-3">
                                            {config.max_concurrent_session}
                                        </div>
                                    </li>
                                    <li className="flex align-items-center px-2 border-300 flex-wrap">
                                        <div className={`text-500 text-right font-medium`}>
                                            Redirect URL:
                                        </div>
                                        <div className="font-normal pl-3">
                                            {config.redirect_url}
                                        </div>
                                    </li>
                                    <Divider align="left" type="dashed">
                                        <b>URLs</b>
                                    </Divider>
                                    <li className="flex align-items-center px-2 border-300 flex-wrap">
                                        <div className={`text-500 text-right font-medium`}>
                                            Portal URL:
                                        </div>
                                        <div className="font-normal pl-3">
                                            {config.external_portal_url ? config.external_portal_url :
                                                config.domain_names?.auth_domain_name ? `https://${config.domain_names?.auth_domain_name}` : `https://${config.status?.secure_ip_address}`}
                                            <span
                                                onClick={() => {
                                                    copy(config.external_portal_url ? config.external_portal_url :
                                                        config.domain_names?.auth_domain_name ? `https://${config.domain_names?.auth_domain_name}` : `https://${config.status?.secure_ip_address}`)
                                                }}
                                                className='ml-2 text-xs text-blue-500 hover:bg-blue-500 hover:text-white transition-duration-500 cursor-pointer border-1 px-1'>
                                                <i className="pi pi-paperclip text-xs"></i>
                                            </span>
                                        </div>
                                    </li>
                                    <li className="flex align-items-center px-2 border-300 flex-wrap">
                                        <div className={`text-500 text-right font-medium`}>
                                            Portal API URL:
                                        </div>
                                        <div className="font-normal pl-3">
                                            {config.domain_names?.auth_domain_name ? `https://${config.domain_names?.auth_domain_name}/api` : `https://${config.status?.egress_ip_address}/api`}
                                            <span
                                                onClick={() => {
                                                    copy(config.domain_names?.auth_domain_name ? `https://${config.domain_names?.auth_domain_name}/api` : `https://${config.status?.egress_ip_address}/api`)
                                                }}
                                                className='ml-2 text-xs text-blue-500 hover:bg-blue-500 hover:text-white transition-duration-500 cursor-pointer border-1 px-1'>
                                                <i className="pi pi-paperclip text-xs"></i>
                                            </span>
                                        </div>
                                    </li>
                                    <li className="flex align-items-center px-2 border-300 flex-wrap">
                                        <div className={`text-500 text-right font-medium`}>
                                            Operator URL:
                                        </div>
                                        <div className="font-normal pl-3">
                                            {config.domain_names?.operator_domain_name ? `https://${config.domain_names?.operator_domain_name}` : `https://${config.status?.egress_ip_address}`}
                                            <span
                                                onClick={() => {
                                                    copy(config.domain_names?.operator_domain_name ? `https://${config.domain_names?.operator_domain_name}` : `https://${config.status?.egress_ip_address}`)
                                                }}
                                                className='ml-2 text-xs text-blue-500 hover:bg-blue-500 hover:text-white transition-duration-500 cursor-pointer border-1 px-1'>
                                                <i className="pi pi-paperclip text-xs"></i>
                                            </span>
                                        </div>
                                    </li>
                                    <li className="flex align-items-center px-2 border-300 flex-wrap">
                                        <div className={`text-500 text-right font-medium`}>
                                            Operator API URL:
                                        </div>
                                        <div className="font-normal pl-3">
                                            {config.domain_names?.operator_domain_name ? `https://${config.domain_names?.operator_domain_name}/api` : `https://${config.status?.egress_ip_address}/api`}
                                            <span
                                                onClick={() => {
                                                    copy(config.domain_names?.operator_domain_name ? `https://${config.domain_names?.operator_domain_name}/api` : `https://${config.status?.egress_ip_address}/api`)
                                                }}
                                                className='ml-2 text-xs text-blue-500 hover:bg-blue-500 hover:text-white transition-duration-500 cursor-pointer border-1 px-1'>
                                                <i className="pi pi-paperclip text-xs"></i>
                                            </span>
                                        </div>
                                    </li>
                                </ul>
                            </div>
                        </div>
                    </div>
                </div>
                <div className="col-12 md:col-6 lg:col-6">
                    <div className="surface-0 p-3 border-1 border-50 border-round" >
                        <div className="flex justify-content-between m1-3">
                            <div>
                                <span className="block text-500 font-medium mb-3">Features</span>
                            </div>
                        </div>
                        <Divider />
                        <div className='flex flex-wrap align-items-center justify-content-center card-container'>
                            <div className="w-full border-round font-bold p-3 ">
                                <ul className="p-0 m-0">
                                    <li className="flex align-items-center py-1 px-2 border-300 flex-wrap">
                                        <div className={`text-500 w-6 md:w-2 font-medium`}>
                                            <i className={`pi pi-circle-fill ${config.ddos_prevention ? "text-green-500" : "text-gray-500"}`}></i>
                                        </div>
                                        <div className="font-normal w-full md:w-8 md:flex-order-0 flex-order-1">DDOS Prevention</div>
                                    </li>
                                    <li className="flex align-items-center text-gray-200 py-1 px-2 border-300 flex-wrap">
                                        <div className={`text-500 w-6 md:w-2 font-medium`}>
                                            <i className={`pi pi-circle-fill text-gray-200`}></i>
                                        </div>
                                        <div className="font-normal w-full md:w-8 md:flex-order-0 flex-order-1">Network Access Control List (future feature)</div>
                                    </li>
                                </ul>
                            </div>
                        </div>
                    </div>
                </div>
            </div>
        </>

    );
}

export default Overview;

