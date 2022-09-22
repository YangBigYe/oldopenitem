// 封装的材料仓库组件
import React, {Component} from "react";
import {Empty, Input, Pagination, Upload} from "antd";
import ModulaCard from "./ModulaCard";
import UpLoadModal from "./UpLoadModal";
import {FileExcelFilled, FileMarkdownFilled, FilePptFilled, FileTextFilled, FileZipFilled, PlusOutlined} from "@ant-design/icons";
import "./MaterialWarehouse.less";

const {Search} = Input;

export default class inedx extends Component {

    state = {
      fileType: [<FileTextFilled key={1} />, <FileExcelFilled key={2} />, <FileZipFilled key={3} />, <FileMarkdownFilled key={4} />, <FilePptFilled key={5} />],
      fileList: [],
      upLoadVisible: false,
      getFileLoading: false,
    }

    componentDidMount() {
      this.getFileList();
    }

    getFileList = () => {
      this.setState({
        getFileLoading: true,
      });
      // request({
      //     url:baseURL+`/review/proj/detailed/${this.props.projectId}`,
      //     // url:`http://49.232.73.36:8081/review/proj/detailed/${this.props.projectId}`,
      //     method:"GET"
      // }).then(res => {
      //     request({
      //         url:baseURL+"/review/query/file",
      //         // url:"http://49.232.73.36:8081/review/query/file",
      //         method:"POST",
      //         data:{
      //             id_list:res.data.materials.files
      //         }
      //     }).then(res => {
      //         this.setState({
      //             fileList:Object.values(res.data),
      //             getFileLoading:false
      //         });
      //     }).catch(err => {
      this.setState({
        getFileLoading: false,
      });
      //     });
      // }).catch(err => {
      //     this.setState({
      //         getFileLoading:false
      //     });
      // });
    }

    fileViewLoader = () => {
      if(this.props.role === "1" && this.state.fileList.length === 0) {
        return (
          <div className="empty-state-box">
            <Empty image={Empty.PRESENTED_IMAGE_SIMPLE} />
          </div>
        );
      }else if((this.props.role === "1" && this.state.fileList.length !== 0) || this.props.role === "1" || this.props.role === "1") {
        return this.state.fileList.map(item => (
          <div className="file-item" key={item.Id} onClick={this.downLoadFile.bind(this, item.uuid)}>
            <div className="icon">
              <FileTextFilled key={item.Id} />
            </div>
            <div className="name">
              <span>{item.name}</span>
            </div>
          </div>
        ));
      }
      return (
        <></>
      );
    }

    render() {
      return (
        <ModulaCard title="材料仓库" right={<Search placeholder="input search text" size="small" style={{width: 200}} />}>
          <div className="material-warehouse-box" data-component="material-warehouse-box">
            <div className="container">
              {
                // this.props.role === "3" || (this.props.role === "4" && this.props.stepName !== "测试框架与论证报告") ? (
                // this.props.stepName == "测试框架与论证报告" ? (
                this.props.role ? (
                  <div className="upload-download-box" onClick={() => {
                    this.setState({
                      upLoadVisible: true,
                    });
                  }}>
                    <Upload
                      name="avatar"
                      listType="picture-card"
                      showUploadList={false}
                      action="https://www.mocky.io/v2/5cc8019d300000980a055e76"
                      openFileDialogOnClick={false}
                    >
                      <div className="file-load-btn" onClick={() => {
                        this.setState({
                          upLoadVisible: true,
                        });
                      }}>
                        <PlusOutlined />
                        <div style={{marginTop: 8, width: "100%"}}>上传</div>
                      </div>
                    </Upload>
                  </div>
                ) : (
                  <></>
                )
              }
              {this.fileViewLoader()}
              <UpLoadModal
                show={this.state.upLoadVisible}
                projectId={this.props.projectId}
                stepId = {this.props.stepId}
                account={this.props.account}
                onClose={() => {
                  this.setState({
                    upLoadVisible: false,
                  });
                }}
                onUpdate = {() => {
                  this.getFileList();
                }}
                getDataList={this.props.getDataList}
              ></UpLoadModal>
            </div>
            <div className="footer">
              <Pagination
                total={85}
                showTotal={total => `Total ${total} items`}
                defaultPageSize={20}
                defaultCurrent={1}
                size="small"
              />
            </div>
          </div>
        </ModulaCard>
      );
    }
}
