import { useEffect, useState } from 'react';
import { useAdminApiConnector } from '../../utils/api-connector';
import { useNavigate } from 'react-router';
import { useToast } from '../../utils/properties';
import VirtualEditor from '../veditor';
import { confirmDialog } from 'primereact/confirmdialog';
import { Dropdown } from 'primereact/dropdown';
import { Button } from 'primereact/button';
import YAML from 'yaml'

const Config = () => {
  /* eslint-disable */
  const [data, setData] = useState("")
  const [mode, setMode] = useState("json")
  const [refresh, setRefresh] = useState(false)
  const navigate = useNavigate()
  const apiInstance = useAdminApiConnector()
  const toast = useToast();

  useEffect(() => {
    apiInstance.api.getConfig()
      .then(res => res.data)
      .then(res => {
        setData(JSON.stringify(res, null, 4))
      })
  }, [])

  const accept = () => {
    let d;
    if (mode == "json") {
      d = JSON.parse(data)
    } else {
      d = YAML.parse(data)
    }

    apiInstance.api.setConfigWithRestartSystem(d)
      .then(() => {
        toast.current.show({ severity: 'success', summary: 'Success', detail: `Updated`, life: 3000 });
        toast.current.show({ severity: 'info', summary: 'Infomation', detail: `Restarting`, life: 3000 });
        setRefresh(!refresh)
      })
      .catch(() => {
        toast.current.show({ severity: 'error', summary: 'Error', detail: `update a config failed`, life: 3000 });
      })
  }

  const reject = () => { }


  return (
    <>
      <div>
        <div className='w-full p-3 bg-gray-50 text-xl font-bold flex justify-content-between'>
          <span className="flex justify-content-start">
            Config Editor
          </span>

          <span className="flex justify-content-end  ">
            <Dropdown value={mode} options={[
              {
                name: "JSON",
                value: "json"
              },
              {
                name: "YAML",
                value: "yaml"
              },
            ]} onChange={(e) => {
              if (mode == "json") {
                const d = JSON.parse(data)
                setData(YAML.stringify(d, null, 4))
              } else {
                const d = YAML.parse(data)
                setData(JSON.stringify(d, null, 4))
              }
              setMode(e.target.value)
            }} optionLabel="name" placeholder="Select a mode" className='mr-2' />
            <Button label="Update" aria-label="Update" onClick={() => {
              confirmDialog({
                message: 'Are you sure you want to update and restart a service?',
                header: 'Confirmation',
                icon: 'pi pi-exclamation-triangle',
                accept,
                reject
              });
            }} />
          </span>
        </div>
        <div className='mb-5'>
          {mode == "json" &&
            <VirtualEditor value={data} onValueChange={(code: string) => {
              setData(code)
            }} lang={"json"} />
          }
          {mode == "yaml" &&
            <VirtualEditor value={data} onValueChange={(code: string) => {
              setData(code)
            }} lang={"yaml"} />
          }
        </div>
      </div>
    </>

  );
}

export default Config;

