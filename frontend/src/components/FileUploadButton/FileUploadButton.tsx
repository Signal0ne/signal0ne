import { ChangeEvent, memo, useState } from 'react';
import { CopyIcon, UploadIcon } from '../Icons/Icons';
import { handleKeyDown } from '../../utils/utils';
import { toast } from 'react-toastify';
import { useUploadWorkflowMutation } from '../../hooks/mutations/useUploadWorkflowMutation';
import ReactModal from 'react-modal';
import yaml, { YAMLException } from 'js-yaml';
import './FileUploadButton.scss';

const customStyles = {
  content: {
    backgroundColor: '#383838',
    borderRadius: '8px',
    bottom: 'auto',
    color: '#fff',
    left: '50%',
    marginRight: '-50%',
    right: 'auto',
    top: '50%',
    transform: 'translate(-50%, -50%)'
  },
  overlay: {
    backgroundColor: 'rgba(0, 0, 0, 0.5)'
  }
};

ReactModal.setAppElement('#root');

const FileUploadButton = () => {
  const [isModalOpen, setIsModalOpen] = useState(false);
  const [webhookUrl, setWebhookUrl] = useState('');

  const openModal = () => setIsModalOpen(true);

  const { mutate } = useUploadWorkflowMutation({
    openModal,
    setWebhookUrl
  });

  const closeModal = () => setIsModalOpen(false);

  const handleFileUpload = (e: ChangeEvent<HTMLInputElement>) => {
    if (e.target.files && e.target.files.length > 0) {
      const file = e.target.files[0];

      if (
        file.type !== 'application/x-yaml' &&
        file.type !== 'application/yaml'
      )
        return;

      const reader = new FileReader();

      reader.onload = async e => {
        try {
          if (!e.target) return;

          const yamlText = e.target.result as string;
          const jsonObject = yaml.load(yamlText) as Record<string, unknown>;

          mutate(jsonObject);
        } catch (error) {
          if (error instanceof YAMLException) {
            toast.error(
              <>
                <p className="toast-title">YAML file parsing error:</p>
                <pre className="toast-code">{error.mark.snippet}</pre>
                <p className="toast-info">
                  Please fix the issue and upload file again
                </p>
              </>,
              {
                autoClose: false,
                className: 'yaml-error-toast'
              }
            );
          }
        }
      };

      reader.readAsText(file);
    }

    e.target.value = '';
  };

  const handleWebhookCopy = async () => {
    try {
      await navigator.clipboard.writeText(webhookUrl);

      toast.success('Webhook URL copied to clipboard');
    } catch (error) {
      toast.error('Failed to copy content to clipboard');
    }
  };

  return (
    <>
      <label className="file-upload-input-container" htmlFor="file-upload">
        <span className="file-upload-label">
          <UploadIcon height={22} width={22} /> Upload Workflow
        </span>
        <input
          accept=".yaml, .yml"
          className="file-upload-input"
          id="file-upload"
          onChange={handleFileUpload}
          type="file"
        />
      </label>
      <ReactModal
        className="webhook-modal-content"
        contentLabel="Your Webhook URL:"
        isOpen={isModalOpen}
        onRequestClose={closeModal}
        style={customStyles}
      >
        <h3 className="modal-title">Your Webhook URL: </h3>
        <div className="modal-content-container">
          <input className="modal-input" readOnly value={webhookUrl} />
          <CopyIcon
            className="modal-copy-icon"
            data-tooltip-class-name="copy-tooltip"
            data-tooltip-content="Copy Webhook URL"
            data-tooltip-id="global"
            height={28}
            onClick={handleWebhookCopy}
            onKeyDown={handleKeyDown(handleWebhookCopy)}
            tabIndex={0}
            width={28}
          />
        </div>
      </ReactModal>
    </>
  );
};

export default memo(FileUploadButton);
