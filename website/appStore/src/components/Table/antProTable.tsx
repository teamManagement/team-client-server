import ProTable, {
  EditableProTable,
  ProTableProps,
  RequestData,
} from "@ant-design/pro-table";
import { Button, FormInstance, Input, message } from "antd";
import {
  CaretUpOutlined,
  CaretDownOutlined,
  SearchOutlined,
} from "@ant-design/icons";
import React, { MutableRefObject, useMemo } from "react";
import classNames from "classnames";
import { Row, Col } from "antd";
import { AntdProTableWrapperOption } from "./types";
import { EditableProTableProps } from "@ant-design/pro-table/lib/components/EditableTable";
import useARH from "use-antd-resizable-header";
import './index.less';

interface AndProTableWrapperStates {
  searchCollapse: boolean;
  searchPlaceholder: string;
}

class AntdProTableWrapper extends React.PureComponent<
  AntdProTableWrapperOption<any, any> & {
    wrapperType: "default" | "editable";
  },
  AndProTableWrapperStates
> {
  /**
   * 查询表单的引用.
   */
  private formRef: MutableRefObject<FormInstance>;

  /**
   * 查询表单的工具栏搜索内容
   */
  private likeVal: string = "";

  public constructor(props: any) {
    super(props);
    this.state = {
      searchCollapse: true,
      searchPlaceholder: "请输入要查询的内容",
    };
    this.formRef = React.createRef() as any;
  }

  /**
   * 重新加载表格
   */
  private reloadTable = () => {
    this.formRef.current.submit();
  };

  /**
   * 高级筛选面板展开切换
   */
  private searchCollapseChange = () => {
    this.setState({
      searchCollapse: !this.state.searchCollapse,
    });
  };

  /**
   * 重置查询表单.
   */
  private restSearchForm = () => {
    this.formRef.current.resetFields();
  };

  /**
   * 高级筛选底部按钮
   * @returns 高级筛选工具
   */
  private searchOptionRender = () => {
    return [
      <Button type="primary" key="query" onClick={this.reloadTable}>
        查询
      </Button>,
      <Button key="rest" onClick={this.restSearchForm}>
        重置
      </Button>,
      // <Button key="cancel" onClick={this.searchCollapseChange}>
      //   取消
      // </Button>,
      <Button
      className="advanced-filter"
      size="small"
      style={{
        height: 13,
        fontSize: 12,
        fontFamily: "Microsoft YaHei",
        fontWeight: 400,
        color: "#0F72EF",
        padding: 0,
        marginRight: 16,
        border: "none",
        lineHeight: "13px",
      }}
      icon={
        this.state.searchCollapse ? (
          <CaretDownOutlined />
        ) : (
          <CaretUpOutlined />
        )
      }
      type="link"
      onClick={this.searchCollapseChange}
    >
      高级筛选
    </Button>
    ];
  };

  /**
   * 请求.
   */
  private request: (params: any) => Promise<RequestData<any>> = async (
    params: any
  ) => {
    this.setState({
      searchCollapse: true,
    });
    params.limit = params.pageSize;
    delete params.pageSize;
    const likeKey = this.props.wrapperSimpleSearchInputName || "likeVal";
    delete params[likeKey];
    if (this.likeVal) {
      params[likeKey] = this.likeVal;
    }
    try {
      if (!this.props.wrapperRequest) {
        message.error("wrapperRequest方法不能为空");
        return {
          success: false,
        } as RequestData<any>;
      }
      const resp = this.props.wrapperRequest(params) as any;
      let respData = resp;
      if (resp instanceof Promise) {
        respData = await resp;
      }

      if (
        !respData ||
        respData instanceof Array ||
        !(respData.data instanceof Array)
      ) {
        if (respData instanceof Array) {
          return {
            success: true,
            data: respData,
          };
        } else {
          return { success: true } as any;
        }
      }

      return {
        success: true,
        data: respData.data,
        total: respData.total,
      } as RequestData<any>;
    } catch (e) {
      message.error((e as any).message);
      return {
        success: false,
      } as RequestData<any>;
    }
  };

  render(): React.ReactNode {
    const className = classNames("custom-pro-table", {
      close: this.state.searchCollapse,
      open: !this.state.searchCollapse,
    });
    return this.props.wrapperSearchDisabled ? (
      this.props.wrapperType === "default" ? (
        <ProTable
          //   {...proTableWrapperProps}
          search={false}
          className="auth-protable"
          request={(this.props as any).dataSource ? undefined : this.request}
          {...this.props}
          formRef={this.formRef}
        />
      ) : (
        <EditableProTable
          //   {...proTableWrapperProps}
          search={false}
          request={this.request}
          {...this.props}
          formRef={this.formRef}
        />
      )
    ) : (
      <div className={className}>
        <Row
          align="middle"
          style={{
            borderRadius: 4,
            height: 50,
            backgroundColor: "#fff",
            paddingLeft: 8,
          }}
        >
          {/* <Col flex="auto">{this.props.wrapperSearchToolbars}</Col> */}
          <Col flex="358px">
            <Row
              align="middle"
              style={{
                height: "100%",
              }}
            >
              {/* <Col flex="auto">
                <Input
                  allowClear
                  style={{
                    width: 260,
                    height: 28,
                    borderRadius: 4,
                    marginRight: "16px",
                  }}
                  placeholder={
                    this.props.wrapperSearchPlaceholder || "请输入要查询的内容"
                  }
                  suffix={<SearchOutlined onClick={this.reloadTable} />}
                  onChange={(e) => (this.likeVal = (e.target as any).value)}
                  onPressEnter={this.reloadTable}
                />
              </Col> */}
              {/* <Col flex="65.7px">
                <Button
                  className="advanced-filter"
                  size="small"
                  style={{
                    height: 13,
                    fontSize: 12,
                    fontFamily: "Microsoft YaHei",
                    fontWeight: 400,
                    color: "#0F72EF",
                    padding: 0,
                    marginRight: 16,
                    border: "none",
                    lineHeight: "13px",
                  }}
                  icon={
                    this.state.searchCollapse ? (
                      <CaretDownOutlined />
                    ) : (
                      <CaretUpOutlined />
                    )
                  }
                  type="link"
                  onClick={this.searchCollapseChange}
                >
                  高级筛选
                </Button>
              </Col> */}
            </Row>
          </Col>
        </Row>
        {this.props.wrapperType === "default" ? (
          <ProTable
            //   {...proTableWrapperProps}
            search={{
              collapsed: this.state.searchCollapse,
              optionRender: this.searchOptionRender,
              labelWidth: "auto",
            }}
            request={(this.props as any).dataSource ? undefined : this.request}
            {...this.props}
            formRef={this.formRef}
          />
        ) : (
          <EditableProTable
            //   {...proTableWrapperProps}
            search={{
              collapsed: this.state.searchCollapse,
              optionRender: this.searchOptionRender,
              labelWidth: "auto",
            }}
            request={this.request}
            {...this.props}
            formRef={this.formRef}
          />
        )}
      </div>
    );
  }
}

/**
 * 使用包装列
 * @param data 数据
 * @returns 包装列
 */
function useWrapperColumns(data: AntdProTableWrapperOption<any, any>) {
  const { components, resizableColumns, tableWidth } = useARH<any>({
    columns: data.wrapperColumns,
    minConstraints: 20,
  });

  return useMemo(() => {
    if (!data.wrapperColumns) {
      return {};
    }
    return { scroll: { x: tableWidth }, components, columns: resizableColumns };
  }, [components, data.wrapperColumns, resizableColumns, tableWidth]);
}

/**
 * 表格包装组件
 * @param props 参数
 * @returns 包装表格
 */
export function AntProTableWrapper<ReqT, DataT>(
  props: AntdProTableWrapperOption<ReqT, DataT> & ProTableProps<any, any>
) {
  const wrapperColumns = useWrapperColumns(props);
  return (
    <AntdProTableWrapper {...props} {...wrapperColumns} wrapperType="default" />
  );
}

/**
 * 可编辑包装表格
 * @param props 参数
 * @returns 包装表格
 */
export function AntEditProTableWrapper<ReqT, DataT>(
  props: AntdProTableWrapperOption<ReqT, DataT> &
    EditableProTableProps<any, any>
) {
  const wrapperColumns = useWrapperColumns(props);
  return (
    <AntdProTableWrapper
      {...props}
      {...wrapperColumns}
      wrapperType="editable"
    />
  );
}
