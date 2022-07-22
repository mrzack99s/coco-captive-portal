import { useEffect, useState } from 'react';
import { useAdminApiConnector } from '../../utils/api-connector';
import { useNavigate } from 'react-router';
import { useToast } from '../../utils/properties';
import { DataTable } from 'primereact/datatable';
import { Column } from 'primereact/column';
import { Button } from 'primereact/button';
import { confirmDialog } from 'primereact/confirmdialog';
import { TypesSessionType } from '../../api';
import { useCookies } from 'react-cookie';
import { FilterMatchMode, FilterOperator } from 'primereact/api';
import { InputText } from 'primereact/inputtext';

const Monitor = () => {
    /* eslint-disable */
    const [data, setData] = useState([] as TypesSessionType[])
    const [refresh, setRefresh] = useState(false)
    const [mounted, setMounted] = useState(false)
    const [loading, setLading] = useState(false)
    const navigate = useNavigate()
    const apiInstance = useAdminApiConnector()
    const toast = useToast();
    const [kickSelected, setKickSelected] = useState({} as TypesSessionType)
    const [cookies, setcookies, removeCookies] = useCookies(["api-token"]);

    const accept = () => {
        apiInstance.api.kickViaIpAddress(kickSelected)
            .then(() => {
                toast.current.show({ severity: 'success', summary: 'Success', detail: `${kickSelected.ip_address} kicked`, life: 3000 });
                setRefresh(!refresh)
            })
            .catch(() => {
                toast.current.show({ severity: 'error', summary: 'Error', detail: `${kickSelected.ip_address} kick failed`, life: 3000 });
            })
    };

    const reject = () => {
        setKickSelected({})
    };

    const kickSession = (e: TypesSessionType) => {
        setKickSelected(e)
    };

    useEffect(() => {
        if (kickSelected.session_uuid) {
            confirmDialog({
                message: 'Do you want to kick this session?',
                header: 'Kick Confirmation',
                icon: 'pi pi-info-circle',
                closable: false,
                accept,
                reject
            });
        }
    }, [kickSelected])

    const getData = () => {
        setLading(true)
        apiInstance.api.getAllSession()
            .then(res => res.data)
            .then(res => {
                setData(res!)
                setLading(false)
            })
    }

    useEffect(() => {
        getData()
    }, [refresh])

    useEffect(() => {
        setInterval(() => {
            getData()
        }, 5000)
    }, [mounted])


    const [filters, setFilters] = useState({
        'global': { value: null, matchMode: FilterMatchMode.CONTAINS },
        'session_uuid': { operator: FilterOperator.OR, constraints: [{ value: null, matchMode: FilterMatchMode.CONTAINS }] },
        'issue': { operator: FilterOperator.OR, constraints: [{ value: null, matchMode: FilterMatchMode.CONTAINS }] },
        'ip_address': { operator: FilterOperator.OR, constraints: [{ value: null, matchMode: FilterMatchMode.CONTAINS }] },
    });

    const [globalFilterValue, setGlobalFilterValue] = useState('');
    const onGlobalFilterChange = (e: { target: { value: any; }; }) => {
        const value = e.target.value;
        let _filters = { ...filters };
        _filters['global'].value = value;

        setFilters(_filters);
        setGlobalFilterValue(value);
    }

    const clearFilter = () => {
        setFilters({
            'global': { value: null, matchMode: FilterMatchMode.CONTAINS },
            'session_uuid': { operator: FilterOperator.AND, constraints: [{ value: null, matchMode: FilterMatchMode.STARTS_WITH }] },
            'issue': { operator: FilterOperator.AND, constraints: [{ value: null, matchMode: FilterMatchMode.STARTS_WITH }] },
            'ip_address': { operator: FilterOperator.AND, constraints: [{ value: null, matchMode: FilterMatchMode.STARTS_WITH }] },
        });
        setGlobalFilterValue('');
    }

    const header = () => {

        return (
            <div className="flex justify-content-between">
                <Button type="button" icon="pi pi-filter-slash" label="Clear" className="p-button-outlined" onClick={clearFilter} />
                <span className="p-input-icon-left">
                    <i className="pi pi-search" />
                    <InputText value={globalFilterValue} onChange={onGlobalFilterChange} placeholder="Keyword Search" />
                </span>
            </div>
        );
    }

    return (
        <div>
            <div>
                <div className='w-full p-3 bg-gray-50 text-xl font-bold'>
                    Session Monitor
                </div>
                <DataTable value={data} header={header}
                    filters={filters}
                    paginator
                    rows={10}
                    globalFilterFields={['session_uuid', 'issue', 'ip_address']}
                    dataKey="session_uuid"
                    responsiveLayout="stack" breakpoint="960px"
                    stateStorage="session" stateKey="dt-state-demo-session" emptyMessage="No session found.">
                    <Column field="session_uuid" header="Session UUID" ></Column>
                    <Column field="issue" header="Issue" ></Column>
                    <Column field="ip_address" header="IP Address"  ></Column>
                    <Column field="last_seen" header="Last Seen"  ></Column>
                    <Column
                        style={{ width: "50px" }}
                        body={(e) => (
                            <>
                                <Button onClick={() => {
                                    kickSession(e)
                                }} icon="pi pi-times" label="Kick" className="p-button-sm p-button-danger p-button-outlined"></Button>
                            </>
                        )}
                    />
                </DataTable>
            </div>
        </div>


    );
}

export default Monitor;

