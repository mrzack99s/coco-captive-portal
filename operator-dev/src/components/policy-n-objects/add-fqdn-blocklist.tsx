import { Button } from "primereact/button";
import { Dialog } from "primereact/dialog"
import { InputText } from "primereact/inputtext";
import { classNames } from "primereact/utils";
import React from "react";
import { FC } from "react";
import { Controller, useForm } from "react-hook-form";
import { useAdminApiConnector } from "../../utils/api-connector";
import { useToast } from "../../utils/properties";

interface props {
    visible: boolean;
    setVisible: React.Dispatch<React.SetStateAction<boolean>>;
    refresh: () => void;
}

export const AddFQDNBlocklist: FC<props> = ({ visible, setVisible, refresh }) => {

    const toast = useToast();
    const defaultValues = {
        fqdn: ""
    };
    const apiInstance = useAdminApiConnector();
    const {
        control,
        formState: { errors, },
        handleSubmit,
        reset,
    } = useForm({ defaultValues });

    const onSubmit = (data: {
        fqdn: string
    }) => {
        add(data.fqdn);
    };


    const add = (fqdn: string) => {
        apiInstance.api.getConfig()
            .then(res => res.data)
            .then(res => {
                let fqdnbl_temp = res.fqdn_blocklist!
                fqdnbl_temp = !fqdnbl_temp ? [] : fqdnbl_temp
                fqdnbl_temp.push(fqdn)
                let conf_temp = res
                conf_temp.fqdn_blocklist = fqdnbl_temp

                apiInstance.api.setConfig(conf_temp)
                    .then(() => {
                        toast.current.show({ severity: 'success', summary: 'Success', detail: `FQDN added`, life: 3000 });
                        reset()
                        setVisible(false)
                        refresh()
                    })
                    .catch(() => {
                        toast.current.show({ severity: 'error', summary: 'Error', detail: `FQDN add failed`, life: 3000 });
                    })
            })

    }

    return (
        <div>
            <Dialog header="Add fqdn to fqdn blocklist" visible={visible} style={{ width: '500px', height: 'auto' }} onHide={() => {
                setVisible(false)
            }}>
                <div className="flex justify-content-center">
                    <div className="card w-10">
                        <form onSubmit={handleSubmit(onSubmit)} className="p-fluid">

                            <div className="field pt-4">
                                <span className="p-float-label">
                                    <Controller name="fqdn" control={control} rules={{ required: 'FQDN is required.' }} render={({ field, fieldState }) => (
                                        <InputText id={field.name} {...field} autoFocus className={classNames({ 'p-invalid': fieldState.invalid })} />
                                    )} />
                                    <label htmlFor="name">FQDN</label>
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