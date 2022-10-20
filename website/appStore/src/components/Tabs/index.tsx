import classNames from "classnames";
import { PureComponent, ReactNode } from "react";
import "./index.less";

export interface TabsHeader {
  key: string;
  title: string;
}
export interface TabsProps {
  currentKey?: string;
  headers: TabsHeader[];
  children: React.ReactNode | React.ReactNode[];
}

interface TabsStates {
  currentKey: string;
}

/**
 * tab切换组件
 */
export class CustomTabs extends PureComponent<TabsProps, TabsStates> {
  constructor(props: TabsProps) {
    super(props);
    this.state = {
      currentKey:
        props.headers && props.headers.length > 0 ? props.headers[0].key : "",
    };
  }

  private headerChange = (key: string) => {
    this.setState({
      currentKey: key,
    });
  };

  render(): ReactNode {
    const currentKey = this.props.currentKey || this.state.currentKey;
    // let children = undefined;
    // if (this.props.children instanceof Array) {
    //   for (let child of this.props.children) {
    //     if ((child as any).key === currentKey) {
    //       children = child;
    //       break;
    //     }
    //   }
    // } else {
    //   const key = (this.props.children as any).key;
    //   if (key === currentKey) {
    //     children = this.props.children;
    //   }
    // }
    const children =
      this.props.children instanceof Array
        ? this.props.children
        : this.props.children
        ? [this.props.children]
        : [];
    return (
      <div className="custom-tabs">
        <div className="custom-tabs-header">
          <ul>
            {this.props.headers.map((h) => (
              <li
                onClick={() => this.headerChange(h.key)}
                key={h.key}
                className={classNames({
                  active: h.key === currentKey,
                })}
              >
                {h.title}
              </li>
            ))}
          </ul>
        </div>
        <div className="custom-tabs-contents">
          {/* <div className="custom-tabs-content">{children}</div> */}
          {children.map((c: any) => (
            <div
              style={{
                display: c.key === currentKey ? "block" : "none",
              }}
              key={c.key}
              className="custom-tabs-content"
            >
              {c}
            </div>
          ))}
        </div>
      </div>
    );
  }
}
