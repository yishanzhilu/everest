interface Roadmap {
  title: string;
  steps: StudyStep[];
}

interface OnlineCourse {
  title: string;
  description: string;
  isOptional?: boolean;
  type: 'online-course';
  options: StudyMaterial[];
}
interface StudyMaterial {
  locales: 'en' | 'zh-CN';
  url: string;
}
interface WebPage {
  title: string;
  description: string;
  isOptional?: boolean;
  type: 'web-page';
  options: StudyMaterial[];
}

interface StudyTopic {
  title: string;
  description: string;
  type: 'parallel-topic' | 'sequence-topic';
  steps: StudyStep[];
}

type StudyStep = OnlineCourse | WebPage | StudyTopic;


// const ml: Roadmap = {
//   title: '机器学习工程师',
//   steps: [
//     {
//       title: '基础课程',
//     },
//     {},
//     {}
//   ]
// }
const roadmap: Roadmap = {
  title: '前端工程师',
  steps: [
    {
      title: 'freeCodeCamp HTML 基础',
      description: '',
      type: 'online-course',
      options: [
        {
          locales: 'en',
          url:
            'https://www.freecodecamp.org/learn/responsive-web-design/basic-html-and-html5/',
        },
        {
          locales: 'zh-CN',
          url:
            'https://learn.freecodecamp.one/responsive-web-design/basic-html-and-html5',
        },
      ],
    },
    {
      title: 'CSS',
      description: '',
      type: 'parallel-topic',
      steps: [
        {
          title: 'freeCodeCamp CSS 基础',
          description: '',
          type: 'web-page',
          options: [
            {
              locales: 'en',
              url:
                'https://www.freecodecamp.org/learn/responsive-web-design/basic-css/',
            },
            {
              locales: 'zh-CN',
              url:
                'https://learn.freecodecamp.one/responsive-web-design/basic-css',
            },
          ],
        },
        {
            title: 'CSS Flex 布局',
            description: '',
            type: 'web-page',
            options: [
              {
                locales: 'en',
                url:
                  'https://css-tricks.com/snippets/css/a-guide-to-flexbox/',
              },
              {
                locales: 'zh-CN',
                url:
                  'http://www.ruanyifeng.com/blog/2015/07/flex-grammar.html',
              },
            ],
          },
      ],
    },
  ],
};

console.log(JSON.stringify(roadmap))