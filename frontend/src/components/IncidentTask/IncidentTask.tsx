import type { IIncidentTask } from '../../contexts/IncidentsProvider/IncidentsProvider';
import type { TaskAssigneeDropdownOption } from '../IncidentPreview/IncidentPreview';
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
import { reorderWithEdge } from '@atlaskit/pragmatic-drag-and-drop-hitbox/util/reorder-with-edge';
import { useCreateIncidentTaskCommentMutation } from '../../hooks/mutations/useCreateIncidentTaskCommentMutation';
import { useIncidentsContext } from '../../hooks/useIncidentsContext';
import { useUpdateIncidentTaskAssigneeMutation } from '../../hooks/mutations/useUpdateIncidentTaskAssignee';
import { useUpdateIncidentTasksPrioritiesMutation } from '../../hooks/mutations/useUpdateIncidentTasksPrioritiesMutation';
import { useUpdateIncidentTaskStatusMutation } from '../../hooks/mutations/useUpdateIncidentTaskStatusMutation';
import AssigneeDropdownOption from '../Dropdown/AssigneeDropdown/AssigneeDropdownOption/AssigneeDropdownOption';
import AssigneeDropdownSingleValue from '../Dropdown/AssigneeDropdown/AssigneeDropdownSingleValue/AssigneeDropdownSingleValue';
import Button from '../Button/Button';
import classNames from 'classnames';
import Dropdown from '../Dropdown/Dropdown';
import Input from '../Input/Input';
import MarkdownWrapper from '../MarkdownWrapper/MarkdownWrapper';
import Spinner from '../Spinner/Spinner';
import TextArea from '../TextArea/TextArea';
import './IncidentTask.scss';

interface IncidentTaskProps {
  availableAssignees: TaskAssigneeDropdownOption[];
  incidentTask: IIncidentTask;
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

  const { selectedIncident } = useIncidentsContext();

  const commentEditorTitleInputRef = useRef<HTMLInputElement>(null);
  const incidentTaskDragHandleRef = useRef(null);
  const incidentTaskRef = useRef(null);

  const handleCloseCommentEditor = () => {
    setCommentTitle('');
    setCommentContent('');
    setIsCommentEditorOpen(false);
  };

  const {
    isPending: isCreateTaskCommentLoading,
    mutate: createTaskCommentMutate
  } = useCreateIncidentTaskCommentMutation({ handleCloseCommentEditor });

  const { mutate: updateTaskAssigneeMutate } =
    useUpdateIncidentTaskAssigneeMutation();

  const { mutate: updateTasksPrioritiesMutate } =
    useUpdateIncidentTasksPrioritiesMutation();

  const { mutate: updateTaskStatusMutate } =
    useUpdateIncidentTaskStatusMutation();

  useEffect(() => {
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

          await updateTasksPrioritiesMutate({
            incidentId: selectedIncident?.id,
            newTasks: {
              incidentTasks: formattedTasksArray
            }
          });
        }
      })
    );
  }, [incidentTask, selectedIncident, updateTasksPrioritiesMutate]);

  useEffect(() => {
    if (isCommentEditorOpen) commentEditorTitleInputRef.current?.focus();
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

  const handleCommentContentChange = (
    e: React.ChangeEvent<HTMLTextAreaElement>
  ) => setCommentContent(e.target.value);

  const handleCommentTitleChange = (e: React.ChangeEvent<HTMLInputElement>) =>
    setCommentTitle(e.target.value);

  const handleOpenCommentEditor = () => setIsCommentEditorOpen(true);

  const handleSaveComment = (incidentItem: IIncidentTask) => {
    if (!commentTitle || !commentContent) return;

    const payload = {
      content: commentContent,
      title: commentTitle
    };

    createTaskCommentMutate({
      incidentId: selectedIncident?.id,
      incidentTaskId: incidentItem.id,
      newComment: payload
    });
  };

  const handleStatusChange = (e: ChangeEvent<HTMLInputElement>) => {
    const updatedTaskStatus = e.target.checked;

    updateTaskStatusMutate({
      incidentId: selectedIncident?.id,
      incidentTaskId: incidentTask.id,
      newStatus: { updatedTaskStatus }
    });
  };

  const handleToggleOpen = () => setIsOpen(prev => !prev);

  const updateTaskAssignee = (
    dropdownOption: TaskAssigneeDropdownOption | null
  ) => {
    if (incidentTask.assignee?.id === dropdownOption?.value.id) return;

    updateTaskAssigneeMutate({
      incidentId: selectedIncident?.id,
      incidentTaskId: incidentTask.id,
      newAssignee: { assignee: dropdownOption?.value }
    });
  };

  const isSaveCommentBtnDisabled =
    !commentTitle || !commentContent || isCreateTaskCommentLoading;

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
              id={incidentTask.taskName}
              onChange={handleStatusChange}
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
            {incidentTask.comments?.length > 0 && (
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
                  className="save-comment-btn"
                  disabled={isSaveCommentBtnDisabled}
                  onClick={() => handleSaveComment(incidentTask)}
                >
                  {isCreateTaskCommentLoading ? <Spinner /> : 'Save Comment'}
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
