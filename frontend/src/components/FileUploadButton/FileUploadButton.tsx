import { ChangeEvent, memo, useEffect, useState } from 'react';
import { CopyIcon, UploadIcon } from '../Icons/Icons';
import { handleKeyDown } from '../../utils/utils';
import { toast } from 'react-toastify';
import { useAuthContext } from '../../hooks/useAuthContext';
import { useWorkflowsContext } from '../../hooks/useWorkflowsContext';
import { Workflow } from '../../data/dummyWorkflows';
import ReactModal from 'react-modal';
import yaml, { YAMLException } from 'js-yaml';
import './FileUploadButton.scss';

interface FetchWorkflowsResponse {
  workflows: Workflow[];
}
interface WorkflowCreateResponse {
  webhook: string;
  workflow: Workflow;
}

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
    backgroundColor: 'rgba(255, 255, 255, 0.6)'
  }
};

ReactModal.setAppElement('#root');

const FileUploadButton = () => {
  const [isModalOpen, setIsModalOpen] = useState(false);
  const [jsonData, setJsonData] = useState<Record<string, unknown> | null>(
    null
  );
  const [webhookUrl, setWebhookUrl] = useState('');

  const { namespaceId } = useAuthContext();
  const { activeWorkflow, setActiveStep, setActiveWorkflow, setWorkflows } =
    useWorkflowsContext();

  useEffect(() => {
    activeWorkflow && setActiveStep(activeWorkflow?.steps[0]);
  }, [activeWorkflow, setActiveStep]);

  const closeModal = () => setIsModalOpen(false);

  const fetchWorkflows = async () => {
    const res = await fetch(
      `${import.meta.env.VITE_SERVER_API_URL}/${namespaceId}/workflow/workflows`
    );

    const data: FetchWorkflowsResponse = await res.json();

    setWorkflows(data.workflows);
  };

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

          if (!namespaceId) throw new Error('Namespace ID not found');

          const yamlText = e.target.result as string;
          const jsonObject = yaml.load(yamlText) as Record<string, unknown>;

          setJsonData(jsonObject);

          const res = await fetch(
            `${import.meta.env.VITE_SERVER_API_URL}/${namespaceId}/workflow/create`,
            {
              body: JSON.stringify(jsonObject),
              headers: {
                'Content-Type': 'application/json'
              },
              method: 'POST'
            }
          );

          const data: WorkflowCreateResponse = await res.json();

          setWebhookUrl(data.webhook);
          setActiveWorkflow(data.workflow);
          openModal();

          await fetchWorkflows();
        } catch (err: unknown) {
          if (err instanceof YAMLException) {
            toast.error(
              <>
                <p className="toast-title">YAML file parsing error:</p>
                <pre className="toast-code">{err.mark.snippet}</pre>
                <p className="toast-info">
                  Please fix the issue and upload file again
                </p>
              </>,
              {
                autoClose: false,
                className: 'yaml-error-toast'
              }
            );

            jsonData && setJsonData(null);
          } else {
            toast.error('Something went wrong.');
          }
        }
      };

      reader.readAsText(file);
    }

    e.target.value = '';
  };

  const handleWebhookCopy = () => {
    navigator.clipboard.writeText(webhookUrl);
    toast.success('Webhook URL copied to clipboard');
  };

  const openModal = () => setIsModalOpen(true);

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
