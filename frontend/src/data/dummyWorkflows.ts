export interface Step {
  condition?: string;
  function: string;
  input: {
    [key: string]: string;
  };
  integration: string;
  name: string;
  output?: {
    [key: string]: string;
  };
}

export interface Workflow {
  description: string;
  id: string;
  name: string;
  steps: Step[];
}

export const DUMMY_WORKFLOWS: Workflow[] = [
  {
    description: 'Lorem ipsum dolor sit amet consectetur adipiscing elit',
    id: '1',
    name: 'Workflow 1',
    steps: [
      {
        name: 'Get ownership',
        integration: 'backstage',
        function: 'get_properties_values',
        input: {
          filter: 'kind=component,metadata.name={service}'
        },
        output: {
          owner: 'spec.owner',
          slack_channel: 'metadata.labels.slack_channel',
          slack_channel2: 'metadata.labels.slack_channel',
          slack_channel3: 'metadata.labels.slack_channel'
        }
      },
      {
        name: 'Inspect traces',
        integration: 'jaeger',
        function: 'log_occurrences',
        input: {
          service: '{service}',
          tags: "{'error': true}"
        },
        output: {
          count: 'sum()',
          log: 'logs:exception.stacktrace | tags:grpc.error_message'
        },
        condition: 'empty(additional_context.log_occurrences)'
      },
      {
        name: 'Send to slack channel',
        integration: 'slack',
        function: 'post_message',
        input: {
          parsable_context_object: '.',
          ignore_context_keys: 'additional_context.components[].slack_channel'
        },
        condition: 'greater(additional_context.log_occurrences.count'
      },
      {
        name: 'Send to slack channel2',
        integration: 'slack',
        function: 'post_message2',
        input: {
          parsable_context_object: '.',
          ignore_context_keys: 'additional_context.components[].slack_channel'
        },
        condition: 'greater(additional_context.log_occurrences.count'
      }
    ]
  },
  {
    description: 'Lorem ipsum dolor sit amet consectetur adipiscing elit',
    id: '2',
    name: 'Workflow 2',
    steps: [
      {
        name: 'Get ownership',
        integration: 'backstage',
        function: 'get_properties_values',
        input: {
          filter: 'kind=component,metadata.name={service}'
        },
        output: {
          owner: 'spec.owner',
          slack_channel: 'metadata.labels.slack_channel'
        }
      },
      {
        name: 'Inspect traces',
        integration: 'jaeger',
        function: 'log_occurrences',
        input: {
          service: '{service}',
          tags: "{'error': true}"
        },
        output: {
          count: 'sum()',
          log: 'logs:exception.stacktrace | tags:grpc.error_message'
        },
        condition: 'empty(additional_context.log_occurrences)'
      },
      {
        name: 'Send to slack channel',
        integration: 'slack',
        function: 'post_message',
        input: {
          parsable_context_object: '.',
          ignore_context_keys: 'additional_context.components[].slack_channel'
        },
        condition: 'greater(additional_context.log_occurrences.count'
      }
    ]
  },
  {
    description: 'Lorem ipsum dolor sit amet consectetur adipiscing elit',
    id: '3',
    name: 'Workflow 3',
    steps: [
      {
        name: 'Get ownership',
        integration: 'backstage',
        function: 'get_properties_values',
        input: {
          filter: 'kind=component,metadata.name={service}'
        },
        output: {
          owner: 'spec.owner',
          slack_channel: 'metadata.labels.slack_channel'
        }
      },
      {
        name: 'Inspect traces',
        integration: 'jaeger',
        function: 'log_occurrences',
        input: {
          service: '{service}',
          tags: "{'error': true}"
        },
        output: {
          count: 'sum()',
          log: 'logs:exception.stacktrace | tags:grpc.error_message'
        },
        condition: 'empty(additional_context.log_occurrences)'
      },
      {
        name: 'Send to slack channel',
        integration: 'slack',
        function: 'post_message',
        input: {
          parsable_context_object: '.',
          ignore_context_keys: 'additional_context.components[].slack_channel'
        },
        condition: 'greater(additional_context.log_occurrences.count'
      }
    ]
  },
  {
    description: 'Lorem ipsum dolor sit amet consectetur adipiscing elit',
    id: '4',
    name: 'Workflow 4',
    steps: [
      {
        name: 'Get ownership',
        integration: 'backstage',
        function: 'get_properties_values',
        input: {
          filter: 'kind=component,metadata.name={service}'
        },
        output: {
          owner: 'spec.owner',
          slack_channel: 'metadata.labels.slack_channel'
        }
      },
      {
        name: 'Inspect traces',
        integration: 'jaeger',
        function: 'log_occurrences',
        input: {
          service: '{service}',
          tags: "{'error': true}"
        },
        output: {
          count: 'sum()',
          log: 'logs:exception.stacktrace | tags:grpc.error_message'
        },
        condition: 'empty(additional_context.log_occurrences)'
      },
      {
        name: 'Send to slack channel',
        integration: 'slack',
        function: 'post_message',
        input: {
          parsable_context_object: '.',
          ignore_context_keys: 'additional_context.components[].slack_channel'
        },
        condition: 'greater(additional_context.log_occurrences.count'
      }
    ]
  },
  {
    description: 'Lorem ipsum dolor sit amet consectetur adipiscing elit',
    id: '5',
    name: 'Workflow 5',
    steps: [
      {
        name: 'Get ownership',
        integration: 'backstage',
        function: 'get_properties_values',
        input: {
          filter: 'kind=component,metadata.name={service}'
        },
        output: {
          owner: 'spec.owner',
          slack_channel: 'metadata.labels.slack_channel'
        }
      },
      {
        name: 'Inspect traces',
        integration: 'jaeger',
        function: 'log_occurrences',
        input: {
          service: '{service}',
          tags: "{'error': true}"
        },
        output: {
          count: 'sum()',
          log: 'logs:exception.stacktrace | tags:grpc.error_message'
        },
        condition: 'empty(additional_context.log_occurrences)'
      },
      {
        name: 'Send to slack channel',
        integration: 'slack',
        function: 'post_message',
        input: {
          parsable_context_object: '.',
          ignore_context_keys: 'additional_context.components[].slack_channel'
        },
        condition: 'greater(additional_context.log_occurrences.count'
      }
    ]
  },
  {
    description: 'Lorem ipsum dolor sit amet consectetur adipiscing elit',
    id: '6',
    name: 'Workflow 6',
    steps: [
      {
        name: 'Get ownership',
        integration: 'backstage',
        function: 'get_properties_values',
        input: {
          filter: 'kind=component,metadata.name={service}'
        },
        output: {
          owner: 'spec.owner',
          slack_channel: 'metadata.labels.slack_channel'
        }
      },
      {
        name: 'Inspect traces',
        integration: 'jaeger',
        function: 'log_occurrences',
        input: {
          service: '{service}',
          tags: "{'error': true}"
        },
        output: {
          count: 'sum()',
          log: 'logs:exception.stacktrace | tags:grpc.error_message'
        },
        condition: 'empty(additional_context.log_occurrences)'
      },
      {
        name: 'Send to slack channel',
        integration: 'slack',
        function: 'post_message',
        input: {
          parsable_context_object: '.',
          ignore_context_keys: 'additional_context.components[].slack_channel'
        },
        condition: 'greater(additional_context.log_occurrences.count'
      }
    ]
  },
  {
    description: 'Lorem ipsum dolor sit amet consectetur adipiscing elit',
    id: '7',
    name: 'Workflow 7',
    steps: [
      {
        name: 'Get ownership',
        integration: 'backstage',
        function: 'get_properties_values',
        input: {
          filter: 'kind=component,metadata.name={service}'
        },
        output: {
          owner: 'spec.owner',
          slack_channel: 'metadata.labels.slack_channel'
        }
      },
      {
        name: 'Inspect traces',
        integration: 'jaeger',
        function: 'log_occurrences',
        input: {
          service: '{service}',
          tags: "{'error': true}"
        },
        output: {
          count: 'sum()',
          log: 'logs:exception.stacktrace | tags:grpc.error_message'
        },
        condition: 'empty(additional_context.log_occurrences)'
      },
      {
        name: 'Send to slack channel',
        integration: 'slack',
        function: 'post_message',
        input: {
          parsable_context_object: '.',
          ignore_context_keys: 'additional_context.components[].slack_channel'
        },
        condition: 'greater(additional_context.log_occurrences.count'
      }
    ]
  }
];
