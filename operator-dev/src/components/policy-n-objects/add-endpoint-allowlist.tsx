import { Button } from "primereact/button";
import { Dialog } from "primereact/dialog"
import { InputText } from "primereact/inputtext";
import { classNames } from "primereact/utils";
import { FC } from "react";
import { Controller, useForm } from "react-hook-form";
import { useAdminApiConnector } from "../../utils/api-connector";
import { useToast } from "../../utils/properties";
interface props {
    visible: boolean;
    setVisible: React.Dispatch<React.SetStateAction<boolean>>;
    refresh: () => void;
}

export const AddEndpointAllowlist: FC<props> = ({ visible, setVisible, refresh }) => {

    const toast = useToast();
    const defaultValues = {
        hostname: "",
        port: ""
    };
    const apiInstance = useAdminApiConnector();
    const {
        control,
        handleSubmit,
        reset,
    } = useForm({ defaultValues });

    const onSubmit = (data: { hostname: string; port: string; }) => {
        add(data.hostname, data.port);
    };


    const add = (hostname: string, port: string) => {
        apiInstance.api.getConfig()
            .then(res => res.data)
            .then(res => {
                let temp = res.allow_endpoints!
                temp = !temp ? [] : temp
                temp.push({
                    hostname: hostname,
                    port: parseInt(port, 10)
                })
                let conf_temp = res
                conf_temp.allow_endpoints = temp

                apiInstance.api.setConfig(conf_temp)
                    .then(() => {
                        toast.current.show({ severity: 'success', summary: 'Success', detail: `Endpoint added`, life: 3000 });
                        reset()
                        setVisible(false)
                        refresh()
                    })
                    .catch(() => {
                        toast.current.show({ severity: 'error', summary: 'Error', detail: `Endpoint add failed`, life: 3000 });
                    })
            })

    }

    return (
        <div>
            <Dialog header="Add endpoint to allow list" visible={visible} style={{ width: '500px', height: 'auto' }} onHide={() => {
                setVisible(false)
            }}>
                <div className="flex justify-content-center">
                    <div className="card w-10">
                        <form onSubmit={handleSubmit(onSubmit)} className="p-fluid">

                            <div className="field pt-4">
                                <span className="p-float-label">
                                    <Controller name="hostname" control={control} rules={{ required: 'Hostname is required.' }} render={({ field, fieldState }) => (
                                        <InputText id={field.name} {...field} autoFocus className={classNames({ 'p-invalid': fieldState.invalid })} />
                                    )} />
                                    <label htmlFor="name">Hostname</label>
                                </span>
                                <span className="p-float-label mt-4">
                                    <Controller name="port" control={control} rules={{ required: 'Port is required.' }} render={({ field, fieldState }) => (
                                        <InputText {...field} className={classNames({ 'p-invalid': fieldState.invalid })} />
                                    )} />
                                    <label htmlFor="name">Port</label>
                                </span>
                            </div>

                            <Button type="submit" label="Add" className="mt-3" />
                        </form>
                    </div>
                </div>
            </Dialog>
        </div>
    )
}