export interface IWorkflowStep {
  condition?: string;
  function: string;
  input: {
    [key: string]: string;
  };
  integration: string;
  integrationType: string;
  name: string;
  displayName: string;
  output?: {
    [key: string]: string;
  };
}

export type IWorkflowTrigger =
  | {
    webhook: {
      output: Record<string, string>;
      condition?: string;
    };
  }
  | {
    scheduled: {
      interval: string;
      output: Record<string, string>;
    };
  };

export interface Workflow {
  description: string;
  id: string;
  name: string;
  steps: IWorkflowStep[];
  trigger: IWorkflowTrigger;
}

export const DUMMY_WORKFLOWS: Workflow[] = [
  {
    description: 'Lorem ipsum dolor sit amet consectetur adipiscing elit',
    id: '1',
    name: 'Workflow 1',
    steps: [
      {
        displayName: 'Get ownership',
        name: 'get-ownership',
        integration: 'backstage',
        integrationType: 'backstage',
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
        displayName: 'Inspect traces',
        name: 'inspect-traces',
        integration: 'jaeger',
        integrationType: 'jaeger',
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
        displayName: 'Send to slack channel',
        name: 'send-slack-message',
        integration: 'slack',
        integrationType: 'slack',
        function: 'post_message',
        input: {
          parsable_context_object: '.',
          ignore_context_keys: 'additional_context.components[].slack_channel'
        },
        condition: 'greater(additional_context.log_occurrences.count'
      }
    ],
    trigger: {
      webhook: {
        output: {
          service: 'job',
          span: 'span_name',
          timestamp: 'startsAt'
        }
      }
    }
  },
  {
    description: 'Lorem ipsum dolor sit amet consectetur adipiscing elit',
    id: '2',
    name: 'Workflow 2',
    steps: [
      {
        displayName: 'Get ownership',
        name: 'get-ownership',
        integration: 'backstage',
        integrationType: 'backstage',
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
        displayName: 'Inspect traces',
        name: 'inspect-traces',
        integration: 'jaeger',
        integrationType: 'jaeger',
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
        displayName: 'Send to slack channel',
        name: 'send-slack-message',
        integration: 'slack',
        integrationType: 'slack',
        function: 'post_message',
        input: {
          parsable_context_object: '.',
          ignore_context_keys: 'additional_context.components[].slack_channel'
        },
        condition: 'greater(additional_context.log_occurrences.count'
      }
    ],
    trigger: {
      scheduled: {
        interval: '15m',
        output: {
          service: 'job',
          span: 'span_name',
          timestamp: 'startsAt'
        }
      }
    }
  },
  {
    description: 'Lorem ipsum dolor sit amet consectetur adipiscing elit',
    id: '3',
    name: 'Workflow 3',
    steps: [
      {
        displayName: 'Get ownership',
        name: 'get-ownership',
        integration: 'backstage',
        integrationType: 'backstage',
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
        displayName: 'Inspect traces',
        name: 'inspect-traces',
        integration: 'jaeger',
        integrationType: 'jaeger',
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
        displayName: 'Send to slack channel',
        name: 'send-slack-message',
        integration: 'slack',
        integrationType: 'slack',
        function: 'post_message',
        input: {
          parsable_context_object: '.',
          ignore_context_keys: 'additional_context.components[].slack_channel'
        },
        condition: 'greater(additional_context.log_occurrences.count'
      }
    ],
    trigger: {
      webhook: {
        output: {
          service: 'job',
          span: 'span_name',
          timestamp: 'startsAt'
        }
      }
    }
  },
  {
    description: 'Lorem ipsum dolor sit amet consectetur adipiscing elit',
    id: '4',
    name: 'Workflow 4',
    steps: [
      {
        displayName: 'Get ownership',
        name: 'get-ownership',
        integration: 'backstage',
        integrationType: 'backstage',
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
        displayName: 'Inspect traces',
        name: 'inspect-traces',
        integration: 'jaeger',
        integrationType: 'jaeger',
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
        displayName: 'Send to slack channel',
        name: 'send-slack-message',
        integration: 'slack',
        integrationType: 'slack',
        function: 'post_message',
        input: {
          parsable_context_object: '.',
          ignore_context_keys: 'additional_context.components[].slack_channel'
        },
        condition: 'greater(additional_context.log_occurrences.count'
      }
    ],
    trigger: {
      webhook: {
        output: {
          service: 'job',
          span: 'span_name',
          timestamp: 'startsAt'
        }
      }
    }
  },
  {
    description: 'Lorem ipsum dolor sit amet consectetur adipiscing elit',
    id: '5',
    name: 'Workflow 5',

    steps: [
      {
        displayName: 'Get ownership',
        name: 'get-ownership',
        integration: 'backstage',
        integrationType: 'backstage',
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
        displayName: 'Inspect traces',
        name: 'inspect-traces',
        integration: 'jaeger',
        integrationType: 'jaeger',
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
        displayName: 'Send to slack channel',
        name: 'send-slack-message',
        integration: 'slack',
        integrationType: 'slack',
        function: 'post_message',
        input: {
          parsable_context_object: '.',
          ignore_context_keys: 'additional_context.components[].slack_channel'
        },
        condition: 'greater(additional_context.log_occurrences.count'
      }
    ],
    trigger: {
      webhook: {
        output: {
          service: 'job',
          span: 'span_name',
          timestamp: 'startsAt'
        }
      }
    }
  },
  {
    description: 'Lorem ipsum dolor sit amet consectetur adipiscing elit',
    id: '6',
    name: 'Workflow 6',

    steps: [
      {
        displayName: 'Get ownership',
        name: 'get-ownership',
        integration: 'backstage',
        integrationType: 'backstage',
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
        displayName: 'Inspect traces',
        name: 'inspect-traces',
        integration: 'jaeger',
        integrationType: 'jaeger',
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
        displayName: 'Send to slack channel',
        name: 'send-slack-message',
        integration: 'slack',
        integrationType: 'slack',
        function: 'post_message',
        input: {
          parsable_context_object: '.',
          ignore_context_keys: 'additional_context.components[].slack_channel'
        },
        condition: 'greater(additional_context.log_occurrences.count'
      }
    ],
    trigger: {
      webhook: {
        output: {
          service: 'job',
          span: 'span_name',
          timestamp: 'startsAt'
        }
      }
    }
  },
  {
    description: 'Lorem ipsum dolor sit amet consectetur adipiscing elit',
    id: '7',
    name: 'Workflow 7',

    steps: [
      {
        displayName: 'Get ownership',
        name: 'get-ownership',
        integration: 'backstage',
        integrationType: 'backstage',
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
        displayName: 'Inspect traces',
        name: 'inspect-traces',
        integration: 'jaeger',
        integrationType: 'jaeger',
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
        displayName: 'Send to slack channel',
        name: 'send-slack-message',
        integration: 'slack',
        integrationType: 'slack',
        function: 'post_message',
        input: {
          parsable_context_object: '.',
          ignore_context_keys: 'additional_context.components[].slack_channel'
        },
        condition: 'greater(additional_context.log_occurrences.count'
      }
    ],
    trigger: {
      webhook: {
        output: {
          service: 'job',
          span: 'span_name',
          timestamp: 'startsAt'
        }
      }
    }
  }
];
