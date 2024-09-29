import { autoScrollForElements } from '@atlaskit/pragmatic-drag-and-drop-auto-scroll/element';
import { FormEvent, useEffect, useRef, useState } from 'react';
import {
  Incident,
  IncidentAssignee
} from '../../contexts/IncidentsProvider/IncidentsProvider';
import { toast } from 'react-toastify';
import { useAuthContext } from '../../hooks/useAuthContext';
import { useIncidentsContext } from '../../hooks/useIncidentsContext';
import AssigneeDropdownOption from '../Dropdown/AssigneeDropdown/AssigneeDropdownOption/AssigneeDropdownOption';
import AssigneeDropdownSingleValueWithImage from '../Dropdown/AssigneeDropdown/AssigneeDropdownSingleValueWithImage/AssigneeDropdownSingleValueWithImage';
import Button from '../Button/Button';
import Dropdown from '../Dropdown/Dropdown';
import IncidentTask from '../IncidentTask/IncidentTask';
import Input from '../Input/Input';
import ReactModal, { Styles } from 'react-modal';
import Spinner from '../Spinner/Spinner';
import './IncidentPreview.scss';

interface IncidentNewTaskResponse {
  updatedIncident: Incident;
}

interface NamespaceUsersResponse {
  users: IncidentAssignee[];
}

export interface TaskAssigneeDropdownOption {
  disabled?: boolean;
  label: string;
  value: IncidentAssignee;
}

const CUSTOM_STYLES: Styles = {
  content: {
    backgroundColor: '#383838',
    border: 'none',
    borderRadius: '8px',
    height: 'max-content',
    margin: 'auto',
    padding: '2rem',
    width: 'max-content'
  },
  overlay: {
    backgroundColor: 'rgba(0, 0, 0, 0.5)',
    zIndex: 9999
  }
};

const IncidentPreview = () => {
  const [availableAssignees, setAvailableAssignees] = useState<
    TaskAssigneeDropdownOption[]
  >([]);
  const [isTaskModalOpen, setIsTaskModalOpen] = useState(false);
  const [taskAssignee, setTaskAssignee] =
    useState<TaskAssigneeDropdownOption | null>(null);
  const [taskErrorMessage, setTaskErrorMessage] = useState('');
  const [taskName, setTaskName] = useState('');

  const previewRef = useRef<HTMLDivElement>(null);

  const { namespaceId } = useAuthContext();
  const { isIncidentPreviewLoading, selectedIncident, setSelectedIncident } =
    useIncidentsContext();

  useEffect(() => {
    const element = previewRef.current;

    if (!element) return;

    return autoScrollForElements({
      element
    });
  }, [selectedIncident]);

  useEffect(() => {
    if (!selectedIncident) return;

    const fetchAvailableAssignees = async () => {
      try {
        const response = await fetch(
          `${import.meta.env.VITE_SERVER_API_URL}/namespace/${namespaceId}/users`
        );

        if (!response.ok) throw new Error('Failed to fetch users');

        const data: NamespaceUsersResponse = await response.json();

        setAvailableAssignees(
          data.users?.map(user => ({
            label: user.name,
            value: user
          }))
        );
      } catch (error) {
        console.error(error);
        setAvailableAssignees([]);
      }
    };

    fetchAvailableAssignees();
  }, [namespaceId, selectedIncident]);

  const handleAddTask = async (e: FormEvent) => {
    e.preventDefault();

    if (!taskAssignee || !taskName) {
      setTaskErrorMessage('Please fill in all the fields');
      return;
    }

    await saveNewTask();
  };

  const handleTaskModalClose = () => {
    setIsTaskModalOpen(false);
    setTaskAssignee(null);
    setTaskName('');
    setTaskErrorMessage('');
  };

  const saveNewTask = async () => {
    if (!taskAssignee || !taskName) return;

    try {
      const response = await fetch(
        `${import.meta.env.VITE_SERVER_API_URL}/${namespaceId}/incident/${selectedIncident?.id}/tasks`,
        {
          body: JSON.stringify({
            assignee: taskAssignee.value,
            comments: [],
            isDone: false,
            items: [],
            priority: selectedIncident?.tasks.length,
            taskName
          }),
          headers: {
            'Content-Type': 'application/json'
          },
          method: 'POST'
        }
      );

      if (!response.ok) throw new Error('Failed to add task');

      const data: IncidentNewTaskResponse = await response.json();

      setSelectedIncident(data.updatedIncident);
      handleTaskModalClose();
      toast.success('Task added successfully');
    } catch (error) {
      if (error instanceof Error) {
        toast.error(error.message);
      } else {
        toast.error('An error occurred while adding the task');
      }
    }
  };

  const getContent = () => {
    if (isIncidentPreviewLoading) return <Spinner />;

    if (!selectedIncident)
      return (
        <p className="incident-preview--empty">
          Please select the incident from the side panel.
        </p>
      );

    return (
      <section className="incident-preview" ref={previewRef}>
        <div className="incident-preview-header">
          <h2 className="incident-preview-header-title">
            {selectedIncident?.title}
          </h2>
          <span className="incident-preview-header-severity">
            <strong>Severity:</strong>{' '}
            <span className="incident-preview-header-severity-value">
              {selectedIncident?.severity}
            </span>
          </span>
          {selectedIncident?.summary && (
            <div className="incident-preview-header-summary">
              <h4 className="incident-preview-header-summary-title">
                Summary:
              </h4>
              <p className="incident-preview-header-summary-content">
                {selectedIncident?.summary}
              </p>
            </div>
          )}
        </div>
        <div className="incident-preview-content">
          <ul className="incident-preview-tasks-list">
            {selectedIncident?.tasks &&
              selectedIncident?.tasks?.map(task => (
                <IncidentTask
                  availableAssignees={availableAssignees}
                  incidentTask={task}
                  key={task.taskName}
                />
              ))}
          </ul>
          {selectedIncident && (
            <Button
              className="add-task-btn"
              onClick={() => setIsTaskModalOpen(true)}
            >
              Add Task
            </Button>
          )}
        </div>
      </section>
    );
  };

  return (
    <main className="incident-preview-container">
      {getContent()}
      <ReactModal
        isOpen={isTaskModalOpen}
        onRequestClose={handleTaskModalClose}
        style={CUSTOM_STYLES}
      >
        <div className="incident-add-task-modal-content">
          <h2 className="incident-task-modal-title">
            Add New Task to <br />
            <span className="incident-title">
              {selectedIncident?.title}
            </span>{' '}
            incident
          </h2>
          <form className="incident-task-form" onSubmit={handleAddTask}>
            <Input
              label="Task Name"
              onChange={e => setTaskName(e.target.value)}
            />
            <Dropdown
              components={{
                Option: AssigneeDropdownOption,
                SingleValue: AssigneeDropdownSingleValueWithImage
              }}
              label="Assignee"
              maxMenuHeight={200}
              menuPortalSelector=".ReactModal__Content"
              menuPosition="fixed"
              onChange={option => option && setTaskAssignee(option)}
              options={availableAssignees}
              placeholder="Select assignee..."
              value={taskAssignee}
            />
            {taskErrorMessage && (
              <p className="error-msg">{taskErrorMessage}</p>
            )}
            <Button disabled={!taskName || !taskAssignee} type="submit">
              Add Task
            </Button>
          </form>
        </div>
      </ReactModal>
    </main>
  );
};

export default IncidentPreview;
