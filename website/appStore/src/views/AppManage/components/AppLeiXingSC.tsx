import { useState, useCallback, useEffect } from 'react'
import { Select } from "antd";

interface ICusFormProps {
    value?: any;
    onChange?: (value?: any) => void;
}

const AppLeiXingSC: React.FC = (props: ICusFormProps) => {

    const [loading, setLoading] = useState<boolean>();
    const [data, setData] = useState<any[]>();

    const fetchData = useCallback(() => {
        setLoading(true);
        //请求数据todo
        setLoading(false);
        setData([{ name: '123', value: '456' }]);
    }, []);

    useEffect(() => {
        fetchData();
        // eslint-disable-next-line react-hooks/exhaustive-deps
    }, [])

    return (
        <Select
            allowClear
            loading={loading}
            style={{ width: 183 }}
            value={props.value}
            onChange={(val) => props.onChange?.(val)}
            options={(data || []).map((m) => { return { label: m.name, value: m.value } })}
            placeholder='应用类型'
        />
    )
};

export default AppLeiXingSC;