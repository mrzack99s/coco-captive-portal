import { Button } from "primereact/button";
import { Dialog } from "primereact/dialog"
import { InputText } from "primereact/inputtext";
import { classNames } from "primereact/utils";
import React from "react";
import { FC } from "react";
import { Controller, useForm } from "react-hook-form";
import { useAdminApiConnector } from "../../utils/api-connector";
import { useToast } from "../../utils/properties";
import { REGEX_NET_CIDR_PATTERN } from "../constants"
interface props {
    visible: boolean;
    setVisible: React.Dispatch<React.SetStateAction<boolean>>;
    refresh: () => void;
}

export const AddUseAuthNetwork: FC<props> = ({ visible, setVisible, refresh }) => {

    const toast = useToast();
    const defaultValues = {
        network: ""
    };
    const apiInstance = useAdminApiConnector();
    const {
        control,
        formState: { errors, },
        handleSubmit,
        reset,
    } = useForm({ defaultValues });

    const onSubmit = (data: {
        network: string
    }) => {
        add(data.network);
    };


    const add = (network: string) => {
        apiInstance.api.getConfig()
            .then(res => res.data)
            .then(res => {
                let allow_network_temp = res.authorized_networks!
                allow_network_temp = !allow_network_temp ? [] : allow_network_temp
                allow_network_temp.push(network)
                let conf_temp = res
                conf_temp.authorized_networks = allow_network_temp

                apiInstance.api.setConfig(conf_temp)
                    .then(() => {
                        toast.current.show({ severity: 'success', summary: 'Success', detail: `network added`, life: 3000 });
                        reset()
                        setVisible(false)
                        refresh()
                    })
                    .catch(() => {
                        toast.current.show({ severity: 'error', summary: 'Error', detail: `network add failed`, life: 3000 });
                    })
            })

    }

    return (
        <div>
            <Dialog header="Add network to use authentication" visible={visible} style={{ width: '530px', height: 'auto' }} onHide={() => {
                setVisible(false)
            }}>
                <div className="flex justify-content-center">
                    <div className="card w-10">
                        <form onSubmit={handleSubmit(onSubmit)} className="p-fluid">

                            <div className="field pt-4">
                                <span className="p-float-label">
                                    <Controller name="network" control={control} rules={{
                                        required: 'Network CIDR is required.',
                                        pattern: { value: REGEX_NET_CIDR_PATTERN, message: 'Invalid Network CIDR' }
                                    }} render={({ field, fieldState }) => (
                                        <InputText id={field.name} {...field} autoFocus className={classNames({ 'p-invalid': fieldState.invalid })} />
                                    )} />
                                    <label htmlFor="name">Network CIDR. Ex 10.0.0.0/24</label>

                                </span>
                                {errors.network &&
                                    <span className="p-error text-xs">{errors.network.message}</span>
                                }

                            </div>

                            <Button type="submit" label="Add" className="mt-3" />
                        </form>
                    </div>
                </div>
            </Dialog>
        </div>
    )
}