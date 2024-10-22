import type { IncidentAssignee } from '../../contexts/IncidentsProvider/IncidentsProvider';
import { autoScrollForElements } from '@atlaskit/pragmatic-drag-and-drop-auto-scroll/element';
import { ChangeEvent, FormEvent, useEffect, useRef, useState } from 'react';
import { toast } from 'react-toastify';
import { useCreateNewIncidentTaskMutation } from '../../hooks/mutations/useCreateNewIncidentTaskMutation';
import { useGetNamespaceUsersQuery } from '../../hooks/queries/useGetNamespaceUsersQuery';
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
  const [isTaskModalOpen, setIsTaskModalOpen] = useState(false);
  const [taskAssignee, setTaskAssignee] =
    useState<TaskAssigneeDropdownOption | null>(null);
  const [taskErrorMessage, setTaskErrorMessage] = useState('');
  const [taskName, setTaskName] = useState('');

  const previewRef = useRef<HTMLDivElement>(null);

  const { isIncidentPreviewLoading, selectedIncident } = useIncidentsContext();

  const { data, isError, isLoading } =
    useGetNamespaceUsersQuery(selectedIncident);

  const handleTaskModalClose = () => {
    setIsTaskModalOpen(false);
    setTaskAssignee(null);
    setTaskName('');
    setTaskErrorMessage('');
  };

  const { isPending: isAddTaskLoading, mutate: saveNewIncidentTaskMutate } =
    useCreateNewIncidentTaskMutation({
      handleTaskModalClose
    });

  useEffect(() => {
    const element = previewRef.current;

    if (!element) return;

    return autoScrollForElements({
      element
    });
  }, [selectedIncident]);

  useEffect(() => {
    if (isError) toast.error('Failed to fetch users');
  }, [isError]);

  const handleAddTask = async (e: FormEvent) => {
    e.preventDefault();

    if (!taskAssignee || !taskName) {
      setTaskErrorMessage('Please fill in all the fields');
      return;
    }

    await saveNewTask();
  };

  const handleTaskNameChange = (e: ChangeEvent<HTMLInputElement>) =>
    setTaskName(e.target.value);

  const openTaskModal = () => setIsTaskModalOpen(true);

  const saveNewTask = async () => {
    if (!selectedIncident || !taskAssignee || !taskName) return;

    const payload = {
      assignee: taskAssignee.value,
      comments: [],
      isDone: false,
      items: [],
      priority: selectedIncident.tasks.length,
      taskName
    };

    await saveNewIncidentTaskMutate({
      incidentId: selectedIncident?.id,
      newIncidentTask: payload
    });
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
                {selectedIncident.summary}
              </p>
            </div>
          )}
        </div>
        <div className="incident-preview-content">
          <ul className="incident-preview-tasks-list">
            {selectedIncident?.tasks &&
              selectedIncident.tasks?.map(task => (
                <IncidentTask
                  availableAssignees={availableAssignees}
                  incidentTask={task}
                  key={task.taskName}
                />
              ))}
          </ul>
          {selectedIncident && (
            <Button className="add-task-btn" onClick={openTaskModal}>
              Add Task
            </Button>
          )}
        </div>
      </section>
    );
  };

  const availableAssignees: TaskAssigneeDropdownOption[] = (
    data?.users ?? []
  ).map(user => ({
    label: user.name,
    value: user
  }));

  const isAddTaskButtonDisabled =
    !taskName || !taskAssignee || isAddTaskLoading;

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
            <Input label="Task Name" onChange={handleTaskNameChange} />
            <Dropdown
              components={{
                Option: AssigneeDropdownOption,
                SingleValue: AssigneeDropdownSingleValueWithImage
              }}
              isDisabled={isLoading}
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
            <Button
              className="add-task-modal-btn"
              disabled={isAddTaskButtonDisabled}
              type="submit"
            >
              {isAddTaskLoading ? <Spinner /> : 'Add Task'}
            </Button>
          </form>
        </div>
      </ReactModal>
    </main>
  );
};

export default IncidentPreview;
