import {
  attachClosestEdge,
  Edge,
  extractClosestEdge
} from '@atlaskit/pragmatic-drag-and-drop-hitbox/closest-edge';
import { ChangeEvent, useEffect, useRef, useState } from 'react';
import { ChevronIcon, UserIcon } from '../Icons/Icons';
import { combine } from '@atlaskit/pragmatic-drag-and-drop/combine';
import {
  draggable,
  dropTargetForElements
} from '@atlaskit/pragmatic-drag-and-drop/element/adapter';
import { DropIndicator } from '@atlaskit/pragmatic-drag-and-drop-react-drop-indicator/box';
import { getIntegrationIcon, handleKeyDown } from '../../utils/utils';
import {
  IIncidentTask,
  Incident
} from '../../contexts/IncidentsProvider/IncidentsProvider';
import { reorderWithEdge } from '@atlaskit/pragmatic-drag-and-drop-hitbox/util/reorder-with-edge';
import { TaskAssigneeDropdownOption } from '../IncidentPreview/IncidentPreview';
import { toast } from 'react-toastify';
import { useAuthContext } from '../../hooks/useAuthContext';
import { useIncidentsContext } from '../../hooks/useIncidentsContext';
import AssigneeDropdownOption from '../Dropdown/AssigneeDropdown/AssigneeDropdownOption/AssigneeDropdownOption';
import AssigneeDropdownSingleValue from '../Dropdown/AssigneeDropdown/AssigneeDropdownSingleValue/AssigneeDropdownSingleValue';
import Button from '../Button/Button';
import classNames from 'classnames';
import Dropdown from '../Dropdown/Dropdown';
import Input from '../Input/Input';
import MarkdownWrapper from '../MarkdownWrapper/MarkdownWrapper';
import TextArea from '../TextArea/TextArea';
import './IncidentTask.scss';

interface AddCommentResponse {
  updatedIncident: Incident;
}

interface IncidentTaskProps {
  availableAssignees: TaskAssigneeDropdownOption[];
  incidentTask: IIncidentTask;
}

interface TaskUpdateResponse {
  updatedIncident: Incident;
}

const IncidentTask = ({
  availableAssignees,
  incidentTask
}: IncidentTaskProps) => {
  const [closestEdge, setClosestEdge] = useState<Edge | null>(null);
  const [commentContent, setCommentContent] = useState('');
  const [commentTitle, setCommentTitle] = useState('');
  const [isCommentEditorOpen, setIsCommentEditorOpen] = useState(false);
  const [isDragging, setIsDragging] = useState(false);
  const [isOpen, setIsOpen] = useState(false);

  const { namespaceId } = useAuthContext();
  const { selectedIncident, setSelectedIncident } = useIncidentsContext();

  const abortControllerRef = useRef(new AbortController());
  const commentEditorTitleInputRef = useRef<HTMLInputElement>(null);

  const incidentTaskDragHandleRef = useRef(null);
  const incidentTaskRef = useRef(null);

  useEffect(() => {
    const handleUpdatePriority = async (newOrderedTasks: IIncidentTask[]) => {
      try {
        const response = await fetch(
          `${import.meta.env.VITE_SERVER_API_URL}/${namespaceId}/incident/${selectedIncident?.id}/update-tasks-priority`,
          {
            body: JSON.stringify({
              incidentTasks: newOrderedTasks
            }),
            headers: {
              'Content-Type': 'application/json'
            },
            method: 'PATCH'
          }
        );

        if (!response.ok) throw new Error('Failed to update the task priority');

        const data: TaskUpdateResponse = await response.json();

        setSelectedIncident(data.updatedIncident);

        toast.success('Task priority updated successfully');
      } catch (error) {
        if (error instanceof Error) {
          toast.error(error.message);
        } else {
          toast.error('An error occurred while updating the task priority');
        }
      }
    };

    const incidentTaskEl = incidentTaskRef.current;
    const incidentTaskDragHandleEl = incidentTaskDragHandleRef.current;

    if (!incidentTaskEl || !incidentTaskDragHandleEl) return;

    return combine(
      draggable({
        element: incidentTaskEl,
        dragHandle: incidentTaskDragHandleEl,
        getInitialData: () => ({
          incidentTaskId: incidentTask.id,
          incidentTask,
          index: incidentTask.priority,
          type: 'taskItem'
        }),
        onDragStart: () => setIsDragging(true),
        onDrop: () => setIsDragging(false)
      }),
      dropTargetForElements({
        element: incidentTaskEl,
        canDrop: ({ source }) => source.data.type === 'taskItem',
        getData: ({ element, input }) => {
          const targetElementId = element.id;
          const foundTaskIndex = selectedIncident?.tasks.findIndex(
            task => task.id === targetElementId
          );

          const data = { index: foundTaskIndex };

          return attachClosestEdge(data, {
            allowedEdges: ['top', 'bottom'],
            element,
            input
          });
        },
        onDrag({ self, source }) {
          const isSource = source.element === incidentTaskEl;

          if (isSource) {
            setClosestEdge(null);
            return;
          }

          const closestEdge = extractClosestEdge(self.data);

          const sourceIndex = source.data.index as number;

          const isItemBeforeSource = incidentTask.priority === sourceIndex - 1;
          const isItemAfterSource = incidentTask.priority === sourceIndex + 1;

          const isDropIndicatorHidden =
            (isItemBeforeSource && closestEdge === 'bottom') ||
            (isItemAfterSource && closestEdge === 'top');

          if (isDropIndicatorHidden) {
            setClosestEdge(null);
            return;
          }

          setClosestEdge(closestEdge);
        },
        onDragLeave() {
          setClosestEdge(null);
        },
        onDrop: async ({ location, self, source }) => {
          setClosestEdge(null);

          if (!selectedIncident) return;

          const closestEdgeOfTarget: Edge | null = extractClosestEdge(
            self.data
          );

          // we shouldn't modify the order if the item is dropped in the same position
          const shouldNotChangeOrder =
            source.data.index === location.current.dropTargets[0].data.index;

          if (shouldNotChangeOrder) return;

          const reorderedArray = reorderWithEdge({
            axis: 'vertical',
            closestEdgeOfTarget: closestEdgeOfTarget,
            indexOfTarget: location.current.dropTargets[0].data.index as number,
            list: selectedIncident.tasks,
            startIndex: source.data.index as number
          });

          const formattedTasksArray = reorderedArray.map((task, index) => ({
            ...task,
            priority: index
          }));

          await handleUpdatePriority(formattedTasksArray);
        }
      })
    );
  }, [incidentTask, namespaceId, selectedIncident, setSelectedIncident]);

  useEffect(() => {
    if (isCommentEditorOpen) {
      commentEditorTitleInputRef.current?.focus();
    }
  }, [isCommentEditorOpen]);

  useEffect(() => {
    if (isDragging) {
      document.body.classList.add('is-dragging');
    } else {
      document.body.classList.remove('is-dragging');
    }
  }, [isDragging]);

  const getAssigneeDropdownValue = () => {
    if (!incidentTask.assignee) return null;

    return {
      disabled: false,
      label: incidentTask.assignee.name,
      value: incidentTask.assignee
    };
  };

  const handleCloseCommentEditor = () => {
    setCommentTitle('');
    setCommentContent('');
    setIsCommentEditorOpen(false);
  };

  const handleCommentContentChange = (
    e: React.ChangeEvent<HTMLTextAreaElement>
  ) => setCommentContent(e.target.value);

  const handleCommentTitleChange = (e: React.ChangeEvent<HTMLInputElement>) =>
    setCommentTitle(e.target.value);

  const handleOpenCommentEditor = () => setIsCommentEditorOpen(true);

  const handleSaveComment = async (incidentItem: IIncidentTask) => {
    if (!commentTitle || !commentContent) return;

    try {
      const response = await fetch(
        `${import.meta.env.VITE_SERVER_API_URL}/${namespaceId}/incident/${selectedIncident?.id}/${incidentItem.id}/add-task-comment`,
        {
          body: JSON.stringify({
            content: commentContent,
            title: commentTitle
          }),
          headers: {
            'Content-Type': 'application/json'
          },
          method: 'POST'
        }
      );

      if (!response.ok) throw new Error('Failed to save the comment');

      const data: AddCommentResponse = await response.json();

      setSelectedIncident(data.updatedIncident);
      setCommentTitle('');
      setCommentContent('');
      setIsCommentEditorOpen(false);
      toast.success('Comment saved successfully');
    } catch (error) {
      if (error instanceof Error) {
        toast.error(error.message);
      } else {
        toast.error('An error occurred while saving the comment');
      }
    }
  };

  const handleStatusChange = async (e: ChangeEvent<HTMLInputElement>) => {
    const updatedTaskStatus = e.target.checked;

    abortControllerRef.current.abort();
    abortControllerRef.current = new AbortController();

    try {
      const response = await fetch(
        `${import.meta.env.VITE_SERVER_API_URL}/${namespaceId}/incident/${selectedIncident?.id}/${incidentTask.id}/status`,
        {
          body: JSON.stringify({
            updatedTaskStatus
          }),
          headers: {
            'Content-Type': 'application/json'
          },
          method: 'PATCH',
          signal: abortControllerRef.current.signal
        }
      );

      if (!response.ok) throw new Error('Failed to update task status');

      const data: TaskUpdateResponse = await response.json();

      setSelectedIncident(data.updatedIncident);
      toast.success('Task status updated successfully');
    } catch (error) {
      if (error instanceof Error) {
        if (error.name !== 'AbortError') toast.error(error.message);
      } else {
        toast.error('An error occurred while updating the task status');
      }
    }
  };

  const handleToggleOpen = () => setIsOpen(prev => !prev);

  const updateTaskAssignee = async (
    dropdownOption: TaskAssigneeDropdownOption | null
  ) => {
    if (incidentTask.assignee?.id === dropdownOption?.value.id) return;

    try {
      const response = await fetch(
        `${import.meta.env.VITE_SERVER_API_URL}/${namespaceId}/incident/${selectedIncident?.id}/${incidentTask.id}/assignee`,
        {
          body: JSON.stringify({
            assignee: dropdownOption?.value
          }),
          headers: {
            'Content-Type': 'application/json'
          },
          method: 'PATCH'
        }
      );

      if (!response.ok) throw new Error('Failed to update task assignee');

      const data: TaskUpdateResponse = await response.json();

      setSelectedIncident(data.updatedIncident);
      toast.success('Task assignee updated successfully');
    } catch (error) {
      if (error instanceof Error) {
        toast.error(error.message);
      } else {
        toast.error('An error occurred while updating the task assignee');
      }
    }
  };

  return (
    <li
      className={classNames('incident-task-container', {
        'is-dragging': isDragging
      })}
      id={incidentTask.id}
      ref={incidentTaskRef}
    >
      <div className="incident-task-tile">
        <div className="incident-task-tile-left">
          <span
            className="incident-task-drag-handle"
            ref={incidentTaskDragHandleRef}
          >
            <span className="dot" />
            <span className="dot" />
            <span className="dot" />
          </span>
          <ChevronIcon
            className={classNames('arrow-icon', {
              open: isOpen
            })}
            height={16}
            onClick={handleToggleOpen}
            onKeyDown={handleKeyDown(handleToggleOpen)}
            tabIndex={0}
            width={16}
          />
          <span className="incident-task-name" onClick={handleToggleOpen}>
            {incidentTask.taskName}
          </span>
        </div>
        <div className="incident-task-tile-right">
          <label
            className={classNames('incident-task-checkbox', {
              checked: incidentTask.isDone
            })}
            htmlFor={incidentTask.taskName}
          >
            <input
              className="incident-task-checkbox-input"
              defaultChecked={incidentTask.isDone}
              onChange={handleStatusChange}
              id={incidentTask.taskName}
              type="checkbox"
            />
            <span className="incident-task-checkbox-label">Done</span>
          </label>
          <Dropdown
            className="assignee-dropdown"
            components={{
              DropdownIndicator: () => null,
              IndicatorsContainer: () => null,
              Option: AssigneeDropdownOption,
              SingleValue: AssigneeDropdownSingleValue
            }}
            isSearchable={false}
            maxMenuHeight={200}
            noOptionsMsg="No users available"
            onChange={option => updateTaskAssignee(option)}
            options={availableAssignees}
            placeholder={<UserIcon height={16} width={16} />}
            value={getAssigneeDropdownValue()}
          />
        </div>
      </div>
      {closestEdge && <DropIndicator edge={closestEdge} gap="16px" />}
      {isOpen && (
        <div className="incident-task-items">
          {incidentTask?.items?.map((incidentItem, index) => (
            <div
              className="incident-task-item"
              key={`${incidentItem.source}-${index}`}
            >
              <div className="incident-task-item-source">
                {getIntegrationIcon(incidentItem.source)}
              </div>
              <div className="incident-task-item-content">
                {incidentItem.content?.map(item => (
                  <div
                    className="incident-task-item-content-info"
                    key={`${item.key}-${item.value}`}
                  >
                    <span className="incident-task-item-content-key">
                      {item.key}:
                    </span>
                    <span className="incident-task-item-content-value">
                      {item.valueType === 'markdown' ? (
                        <MarkdownWrapper content={String(item.value)} />
                      ) : (
                        item.value
                      )}
                    </span>
                  </div>
                ))}
              </div>
            </div>
          ))}
          <div className="users-comments-container">
            <h3 className="users-comments-header-title">Users Comments</h3>
            <hr className="users-comments-separator" />
            {incidentTask.comments.length > 0 && (
              <div className="users-comments-body">
                {incidentTask.comments.map((comment, index) => (
                  <div className="user-comment-container" key={index}>
                    <div
                      className="user-comment-content-info"
                      key={`${comment.content.key}-${comment.content.value}`}
                    >
                      <span className="user-comment-content-key">
                        {comment.content.key}
                      </span>
                      <span className="user-comment-content-value">
                        <MarkdownWrapper
                          content={String(comment.content.value)}
                        />
                      </span>
                    </div>
                    <i className="user-comment-author-info">
                      <span className="user-comment-author-name">
                        ~ {comment.source.name}
                      </span>
                      <span className="user-comment-date">
                        {new Date(comment.timestamp * 1000).toLocaleDateString(
                          'en-US',
                          {
                            day: '2-digit',
                            month: 'short',
                            year: 'numeric',
                            hour: '2-digit',
                            minute: '2-digit'
                          }
                        )}
                      </span>
                    </i>
                  </div>
                ))}
              </div>
            )}
          </div>
          {isCommentEditorOpen && (
            <div className="comment-editor-container">
              <Input
                label="Comment Title"
                onChange={handleCommentTitleChange}
                placeholder="Enter comment title"
                ref={commentEditorTitleInputRef}
                type="text"
                value={commentTitle}
              />
              <TextArea
                label="Comment Content"
                onChange={handleCommentContentChange}
                placeholder="Type your comment here..."
                value={commentContent}
              />
            </div>
          )}
          <div className="comment-editor-buttons">
            {isCommentEditorOpen && (
              <>
                <Button onClick={handleCloseCommentEditor}>
                  Discard Comment
                </Button>
                <Button
                  aria-disabled={!commentTitle || !commentContent}
                  disabled={!commentTitle || !commentContent}
                  onClick={() => handleSaveComment(incidentTask)}
                >
                  Save Comment
                </Button>
              </>
            )}
            {!isCommentEditorOpen && (
              <Button onClick={handleOpenCommentEditor}>Add Comment</Button>
            )}
          </div>
        </div>
      )}
    </li>
  );
};

export default IncidentTask;
