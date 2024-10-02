import Markdown from 'react-markdown';
import './MarkdownWrapper.scss';

interface MarkdownWrapperProps {
  content: string;
}

const MarkdownWrapper = ({ content }: MarkdownWrapperProps) => (
  <Markdown
    children={content}
    className="markdown-container"
    components={{
      // eslint-disable-next-line @typescript-eslint/no-unused-vars
      a: ({ node, children, ...rest }) => (
        <a {...rest} className="markdown-link" rel="noreferrer" target="_blank">
          {children}
        </a>
      ),
      // eslint-disable-next-line @typescript-eslint/no-unused-vars
      code: ({ node, children, ...rest }) => (
        <code {...rest} className="markdown-code">
          {children}
        </code>
      ),
      // eslint-disable-next-line @typescript-eslint/no-unused-vars
      p: ({ node, children, ...props }) => (
        <p {...props} className="markdown-paragraph">
          {children}
        </p>
      )
    }}
  />
);

export default MarkdownWrapper;
