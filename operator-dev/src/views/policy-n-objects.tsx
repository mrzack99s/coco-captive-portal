import Navbar from "../components/navbar";
import { useEffect, useState } from "react";
import { useAdminApiConnector } from "../utils/api-connector";
import { VirtualScroller } from "primereact/virtualscroller";
import { CONSTANT_FQDN_BLOCKLIST, CONSTANT_BYPASS_NETWORK, CONSTANT_ENDPOINT_ALLOWLIST } from "../components/constants";
import { Fieldset } from "primereact/fieldset";
import { Chip } from "primereact/chip";
import { TypesConfigType, TypesEndpointType } from "../api";
import { Button } from "primereact/button";
import { InputText } from "primereact/inputtext";
import { Tooltip } from "primereact/tooltip";
import { AddFQDNBlocklist } from "../components/policy-n-objects/add-fqdn-blocklist";
import { AddEndpointAllowlist } from "../components/policy-n-objects/add-endpoint-allowlist";
import { AddBypassNetwork } from "../components/policy-n-objects/add-bypass-network";
import { useToast } from "../utils/properties";
import { confirmDialog } from "primereact/confirmdialog";
import { Copyright } from "../components/copyright"

const PolicyNObjectView = () => {
    const [config, setConfig] = useState({} as TypesConfigType)
    const apiInstance = useAdminApiConnector()
    const [refresh, setRefresh] = useState(false)
    const [loading, setLoading] = useState(false)
    const [filterEndpointAllowList, setFilterEndpointAllowList] = useState("")
    const [endpointAllowList, setEndpointAllowList] = useState([] as TypesEndpointType[])
    const [filterFQDNBlockList, setFilterFQDNBlockList] = useState("")
    const [fqdnBlockList, setFQDNBlockList] = useState([] as string[])
    const [filterBypassNetwork, setFilterBypassNetwork] = useState("")
    const [bypassNetwork, setBypassNetwork] = useState([] as string[])
    const [dialogBypassNetwork, setDialogBypassNetwork] = useState(false)
    const [dialogFqdnBlocklist, setDialogFqdnBlocklist] = useState(false)
    const [dialogEndpointAllowlist, setDialogEndpointAllowlist] = useState(false)
    const [deleteSelected, setDeleteSelected] = useState("" as string | object)
    const [deleteMode, setDeleteMode] = useState("")
    const toast = useToast();

    useEffect(() => {
        getData()
        // eslint-disable-next-line react-hooks/exhaustive-deps
    }, [])

    const getData = () => {
        setLoading(true)
        apiInstance.api.getConfig()
            .then(res => res.data)
            .then(res => {
                setConfig(res)
                setEndpointAllowList(res.allow_endpoints!)
                setFQDNBlockList(res.fqdn_blocklist!)
                setBypassNetwork(res.bypass_networks!)
                setLoading(false)
            })
    }

    useEffect(() => {
        getData()
        // eslint-disable-next-line react-hooks/exhaustive-deps
    }, [refresh])

    const filter = (kw: string, list: any, setFunc: React.SetStateAction<any>) => {
        if (kw !== "") {

            // eslint-disable-next-line array-callback-return
            const filtered = list.filter((e: any) => {
                if (typeof e === "string") {
                    return e.includes(kw)
                } else {
                    const objKeys = Object.keys(e)
                    let found = false
                    for (let i = 0; i < objKeys.length; i++) {
                        if (typeof e[objKeys[i]] === "string") {
                            if (e[objKeys[i]].includes(kw)) {
                                found = true
                                break
                            }
                        }
                        if (typeof e[objKeys[i]] === "number") {
                            if (e[objKeys[i]].toString().includes(kw)) {
                                found = true
                                break
                            }
                        }
                    }
                    return found
                }
            })
            console.log(filtered)
            setFunc(filtered!)
        }
    }

    const deleteAccept = () => {
        let temp = config;
        let index = -1;
        switch (deleteMode) {
            case CONSTANT_FQDN_BLOCKLIST:
                index = temp.fqdn_blocklist!.findIndex(e => e === deleteSelected)
                temp.fqdn_blocklist!.splice(index, 1)
                break;
            case CONSTANT_ENDPOINT_ALLOWLIST:
                index = temp.allow_endpoints!.findIndex(e => e === deleteSelected)
                temp.allow_endpoints!.splice(index, 1)
                break;
            case CONSTANT_BYPASS_NETWORK:
                index = temp.bypass_networks!.findIndex(e => e === deleteSelected)
                temp.bypass_networks!.splice(index, 1)
                break;
        }

        apiInstance.api.setConfig(temp)
            .then(() => {
                setRefresh(!refresh)
                toast.current.show({ severity: 'success', summary: 'Success', detail: `Deleted`, life: 3000 });
            })
            .catch(() => {
                toast.current.show({ severity: 'error', summary: 'Error', detail: `Delete add failed`, life: 3000 });
            })
    }

    const deleteReject = () => {
        setDeleteMode("")
        setDeleteSelected({})
    }

    useEffect(() => {
        if (filterEndpointAllowList !== "") {
            filter(filterEndpointAllowList, config.allow_endpoints, setEndpointAllowList)
        }
        // eslint-disable-next-line react-hooks/exhaustive-deps
    }, [filterEndpointAllowList])

    useEffect(() => {
        if (filterFQDNBlockList !== "") {
            filter(filterFQDNBlockList, config.fqdn_blocklist, setFQDNBlockList)
        }
        // eslint-disable-next-line react-hooks/exhaustive-deps
    }, [filterFQDNBlockList])

    useEffect(() => {
        if (filterBypassNetwork !== "") {
            filter(filterBypassNetwork, config.bypass_networks, setBypassNetwork)
        }
        // eslint-disable-next-line react-hooks/exhaustive-deps
    }, [filterBypassNetwork])

    useEffect(() => {
        if (deleteSelected && deleteMode) {
            switch (deleteMode) {
                case CONSTANT_FQDN_BLOCKLIST:
                    confirmDialog({
                        message: 'Do you want to delete this fqdn?',
                        header: 'Delete Confirmation',
                        icon: 'pi pi-info-circle',
                        acceptClassName: 'p-button-danger',
                        closable: false,
                        draggable: false,
                        accept: deleteAccept,
                        reject: deleteReject
                    });
                    break;
                case CONSTANT_ENDPOINT_ALLOWLIST:
                    confirmDialog({
                        message: 'Do you want to delete this endpoint?',
                        header: 'Delete Confirmation',
                        icon: 'pi pi-info-circle',
                        acceptClassName: 'p-button-danger',
                        closable: false,
                        draggable: false,
                        accept: deleteAccept,
                        reject: deleteReject
                    });
                    break;
                case CONSTANT_BYPASS_NETWORK:
                    confirmDialog({
                        message: 'Do you want to delete this network address?',
                        header: 'Delete Confirmation',
                        icon: 'pi pi-info-circle',
                        acceptClassName: 'p-button-danger',
                        closable: false,
                        draggable: false,
                        accept: deleteAccept,
                        reject: deleteReject
                    });
                    break;
            }
        }
        // eslint-disable-next-line react-hooks/exhaustive-deps
    }, [deleteSelected, deleteMode])



    const clearFilter = (fMode: string) => {
        switch (fMode) {
            case CONSTANT_FQDN_BLOCKLIST:
                setFilterFQDNBlockList("")
                setFQDNBlockList(config.fqdn_blocklist!)
                break;
            case CONSTANT_ENDPOINT_ALLOWLIST:
                setFilterEndpointAllowList("")
                setEndpointAllowList(config.allow_endpoints!)
                break;
            case CONSTANT_BYPASS_NETWORK:
                setFilterBypassNetwork("")
                setBypassNetwork(config.bypass_networks!)
                break;
        }

    }

    return (
        <div className="mb-5">
            <Navbar />
            <AddFQDNBlocklist refresh={() => {
                setRefresh(!refresh)
            }} visible={dialogFqdnBlocklist} setVisible={setDialogFqdnBlocklist} />
            <AddEndpointAllowlist refresh={() => {
                setRefresh(!refresh)
            }} visible={dialogEndpointAllowlist} setVisible={setDialogEndpointAllowlist} />
            <AddBypassNetwork refresh={() => {
                setRefresh(!refresh)
            }} visible={dialogBypassNetwork} setVisible={setDialogBypassNetwork} />

            <div className="grid grid-nogutter m-0" style={{ position: "relative", top: "65px" }}>
                <div className="col hidden lg:inline grid-nogutter"></div>
                <div className="col-12 lg:col-8 grid-nogutter">
                    <div className='w-full p-3 bg-gray-50 text-xl font-bold'>
                        Policy and Objects
                    </div>

                    <div className="grid nested-grid mt-3 h-36rem px-2 mb-3">
                        <div className="col-12 lg:col-5 lg:inline p-0">
                            <Fieldset legend="Ingress" className="custom-fieldset-p0 border-1 border-gray-50 border-round-2xl h-full">
                                <div className="p-2 font-medium text-xs text-gray-500">
                                    {"Bypass Captive Portal".toUpperCase()}
                                </div>
                                <div className="grid p-2 font-medium text-xs text-gray-500">
                                    <div className="col-2 lg:col-1 p-0 my-auto mx-auto">
                                        <Button tooltip="Add network to bypass a captive portal"
                                            tooltipOptions={{ mouseTrack: true }}
                                            icon="pi pi-plus-circle"
                                            className="p-button-primary w-full py-2"
                                            onClick={() => setDialogBypassNetwork(true)} />
                                    </div>
                                    <div className="col-2 lg:col-1 p-0 my-auto mx-auto pl-1">
                                        <Button tooltip="Filter cleanup"
                                            tooltipOptions={{ mouseTrack: true }}
                                            icon="pi pi-filter-slash"
                                            className="p-button-secondary w-full py-2"
                                            onClick={() => clearFilter(CONSTANT_BYPASS_NETWORK)} />
                                    </div>
                                    <div className="col-8 lg:col-10 p-0 my-auto mx-auto">
                                        <InputText className="w-full" value={filterBypassNetwork} onChange={(e) => { setFilterBypassNetwork(e.target.value) }} placeholder="Keyword Search" />
                                    </div>
                                </div>
                                <div className="p-2 font-medium text-xs text-gray-500 border-1 border-gray-50">
                                    {bypassNetwork &&
                                        <VirtualScroller
                                            loading={loading}
                                            className="h-18rem"
                                            items={bypassNetwork}
                                            itemSize={50} itemTemplate={(item) => (
                                                <div className="grid m-0 border-bottom-1 border-gray-50 py-1">
                                                    <div className="col-10 p-0 lg:p-1">
                                                        <Chip label={item} />
                                                    </div>
                                                    <div className="col-2 p-0 m-auto text-right">
                                                        <Tooltip target=".tooltip-tracking" mouseTrack mouseTrackLeft={10} />
                                                        <span
                                                            onClick={() => {
                                                                setDeleteSelected(item)
                                                                setDeleteMode(CONSTANT_BYPASS_NETWORK)
                                                                confirmDialog({
                                                                    message: 'Do you want to delete this network address?',
                                                                    header: 'Delete Confirmation',
                                                                    icon: 'pi pi-info-circle',
                                                                    acceptClassName: 'p-button-danger',
                                                                    closable: false,
                                                                    draggable: false,
                                                                    accept: deleteAccept,
                                                                    reject: deleteReject
                                                                });
                                                            }}
                                                            className='text-sm text-red-500 hover:bg-red-500 hover:text-white transition-duration-500 cursor-pointer border-1 p-1 tooltip-tracking'
                                                            data-pr-tooltip="Delete"
                                                        >
                                                            <i className="pi pi-trash text-sm"></i>
                                                        </span>
                                                    </div>
                                                </div>
                                            )}
                                            showLoader delay={250} />
                                    }
                                    {!bypassNetwork &&
                                        <div className="flex align-items-center justify-content-center h-18rem">
                                            Not found network
                                        </div>
                                    }

                                </div>
                            </Fieldset>
                        </div>
                        <div className="col-12 hidden lg:col-2 lg:flex align-items-center justify-content-center">
                            <span className="text-center">
                                <p className="m-0">
                                    <i className="pi pi-arrow-right" style={{ 'fontSize': '1em' }}></i>
                                </p>
                                <p className="m-0">
                                    Next
                                </p>
                            </span>
                        </div>
                        <div className="col-12 mt-2 mb-2 lg:mt-0 lg:mb-0 lg:col-5 lg:inline p-0">
                            <Fieldset legend="Egress" className="custom-fieldset-p0 border-1 border-gray-50 border-round-2xl h-full">

                                <div className="p-2 font-medium text-xs text-gray-500">
                                    {"Bypass Endpoint List".toUpperCase()}
                                </div>
                                <div className="grid p-2 font-medium text-xs text-gray-500">
                                    <div className="col-2 lg:col-1 p-0 my-auto mx-auto">
                                        <Button tooltip="Add endpoint to allow list" tooltipOptions={{ mouseTrack: true }}
                                            icon="pi pi-plus-circle" className="p-button-primary w-full py-2"
                                            onClick={() => setDialogEndpointAllowlist(true)} />
                                    </div>
                                    <div className="col-2 lg:col-1 p-0 my-auto mx-auto pl-1">
                                        <Button tooltip="Filter cleanup" tooltipOptions={{ mouseTrack: true }}
                                            icon="pi pi-filter-slash" className="p-button-secondary w-full py-2" onClick={() => clearFilter(CONSTANT_ENDPOINT_ALLOWLIST)} />
                                    </div>
                                    <div className="col-8 lg:col-10 p-0 my-auto mx-auto">
                                        <InputText className="w-full" value={filterEndpointAllowList} onChange={(e) => { setFilterEndpointAllowList(e.target.value) }} placeholder="Keyword Search" />
                                    </div>
                                </div>
                                <div className="p-2 font-medium text-xs text-gray-500 border-1 border-gray-50">
                                    {endpointAllowList &&
                                        <VirtualScroller
                                            loading={loading}
                                            className="h-18rem"
                                            items={endpointAllowList}
                                            itemSize={50} itemTemplate={(item) => (
                                                <div className="grid m-0 border-bottom-1 border-gray-50 py-1">
                                                    <div className="col-10 p-0 lg:p-1">
                                                        <div className="grid grid-nogutter">
                                                            <div className="col-12 lg:col">
                                                                <Chip label={item.hostname} />
                                                            </div>
                                                            <div className="col-12 lg:col mt-1 lg:mt-0">
                                                                <Chip label={item.port} />
                                                            </div>
                                                        </div>
                                                    </div>
                                                    <div className="col-2 p-0 m-auto text-right">
                                                        <Tooltip target=".tooltip-tracking" mouseTrack mouseTrackLeft={10} />
                                                        <span
                                                            onClick={() => {
                                                                setDeleteSelected(item)
                                                                setDeleteMode(CONSTANT_ENDPOINT_ALLOWLIST)
                                                                confirmDialog({
                                                                    message: 'Do you want to delete this endpoint?',
                                                                    header: 'Delete Confirmation',
                                                                    icon: 'pi pi-info-circle',
                                                                    acceptClassName: 'p-button-danger',
                                                                    closable: false,
                                                                    draggable: false,
                                                                    accept: deleteAccept,
                                                                    reject: deleteReject
                                                                });
                                                            }}
                                                            className='text-sm text-red-500 hover:bg-red-500 hover:text-white transition-duration-500 cursor-pointer border-1 p-1 tooltip-tracking'
                                                            data-pr-tooltip="Delete"
                                                        >
                                                            <i className="pi pi-trash text-sm"></i>
                                                        </span>
                                                    </div>
                                                </div>
                                            )}
                                            showLoader delay={250} />
                                    }
                                    {!endpointAllowList &&
                                        <div className="flex align-items-center justify-content-center h-18rem">
                                            Not found endpoint
                                        </div>
                                    }

                                </div>

                                <div className="p-2 font-medium text-xs text-gray-500">
                                    {"FQDN block list".toUpperCase()}
                                </div>
                                <div className="grid p-2 font-medium text-xs text-gray-500">
                                    <div className="col-2 lg:col-1 p-0 my-auto mx-auto">
                                        <Button tooltip="Add fqdn to block list"
                                            tooltipOptions={{ mouseTrack: true }} icon="pi pi-plus-circle"
                                            className="p-button-primary w-full py-2" onClick={() => setDialogFqdnBlocklist(true)} />
                                    </div>
                                    <div className="col-2 lg:col-1 p-0 my-auto mx-auto pl-1">
                                        <Button tooltip="Filter cleanup" tooltipOptions={{ mouseTrack: true }}
                                            icon="pi pi-filter-slash" className="p-button-secondary w-full py-2" onClick={() => clearFilter(CONSTANT_FQDN_BLOCKLIST)} />
                                    </div>
                                    <div className="col-8 lg:col-10 p-0 my-auto mx-auto">
                                        <InputText className="w-full" value={filterFQDNBlockList} onChange={(e) => { setFilterFQDNBlockList(e.target.value) }} placeholder="Keyword Search" />
                                    </div>
                                </div>
                                <div className="p-2 font-medium text-xs text-gray-500 border-1 border-gray-50">
                                    {fqdnBlockList &&
                                        <VirtualScroller
                                            loading={loading}
                                            className="h-18rem"
                                            items={fqdnBlockList}
                                            itemSize={50} itemTemplate={(item) => (
                                                <div className="grid m-0 border-bottom-1 border-gray-50 py-1">
                                                    <div className="col-10 p-0 lg:p-1">
                                                        <Chip label={item} />
                                                    </div>
                                                    <div className="col-2 p-0 m-auto text-right">
                                                        <Tooltip target=".tooltip-tracking" mouseTrack mouseTrackLeft={10} />
                                                        <span
                                                            onClick={() => {
                                                                setDeleteSelected(item)
                                                                setDeleteMode(CONSTANT_FQDN_BLOCKLIST)
                                                                confirmDialog({
                                                                    message: 'Do you want to delete this fqdn?',
                                                                    header: 'Delete Confirmation',
                                                                    icon: 'pi pi-info-circle',
                                                                    acceptClassName: 'p-button-danger',
                                                                    closable: false,
                                                                    draggable: false,
                                                                    accept: deleteAccept,
                                                                    reject: deleteReject
                                                                });
                                                            }}
                                                            className='text-sm text-red-500 hover:bg-red-500 hover:text-white transition-duration-500 cursor-pointer border-1 p-1 tooltip-tracking'
                                                            data-pr-tooltip="Delete"
                                                        >
                                                            <i className="pi pi-trash text-sm"></i>
                                                        </span>
                                                    </div>
                                                </div>
                                            )}
                                            showLoader delay={250} />
                                    }
                                    {!fqdnBlockList &&
                                        <div className="flex align-items-center justify-content-center h-18rem">
                                            Not found endpoint
                                        </div>
                                    }

                                </div>
                            </Fieldset>

                        </div>
                    </div>
                </div>
                <div className="col hidden lg:inline grid-nogutter"></div>
            </div>
            <Copyright />
        </div>
    );
};

export default PolicyNObjectView;