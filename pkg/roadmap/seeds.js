const schema = {
  title: '',
  steps: [
    {
      title: 'step 1',
      type: 'parallel-group | online-course | book | self-study',
      lang: [
        {
          locale: 'en | zh-CN',
          price: '',
        }
      ]

    },
    {
      title: 'step 2',
      type: 'parallel-group',
      subSteps: [],
    },
    {
      title: 'step 1',
      type: 'sequncial-group',
      subSteps: [],
    },
  ],
};

const roadmaps = [
  {
    title: '机器学习算法工程师',
    missions: [
      {
        title: 'Learn basic Python:',
        desc: '',
        type: 'group',
        children: [
          {
            title: 'Python for Everybody',
            type: 'mooc',
            desc:
              'Learn to Program and Analyze Data with Python. Develop programs to gather, clean, analyze, and visualize data.',
            time: '12 小时',
            url: 'https: //www.coursera.org/learn/python',
          },
          {
            type: 'book',
            title: 'Learn python the Hard way',
            desc:
              'A Very Simple Introduction to the Terrifyingly Beautiful World of Computers and Code',
            url: 'https: //learnpythonthehardway.org/book/',
            optional: true,
          },
        ],
      },
      {
        title: 'Learn basic data structures using Python',
        desc: `For this go to Geekforgeeks-Data structures, open each data structure and go to the Python tab. See how they implement it, learn from them.
              
                  DO NOT JUST MINDLESSLY LOOK AT THE CODE, write the code and learn from it.`,
        url: 'https: //www.geeksforgeeks.org/data-structures/',
        type: 'self-study',
        time: '30 小时',
      },
      {
        title: 'Learn some cool Python packages',
        desc: `Do basic Linear Algebra, Statistics and Probability with Python`,
        type: 'group',
        children: [
          {
            title: 'Linear Algebra with Numpy',
            desc: `The Linear Algebra module of NumPy offers various methods to apply linear algebra on any numpy array.
                      One can find:
                      
                      rank, determinant, trace, etc. of an array.
                      eigen values of matrices
                      matrix and vector products (dot, inner, outer,etc. product), matrix exponentiation
                      solve linear or tensor equations and much more!`,
            url:
              'https: //docs.scipy.org/doc/numpy/reference/routines.linalg.html',
            time: '5 小时',
          },
        ],
      },
      {
        title: 'o basic Linear Algebra, Statistics and Probability with Python',
        desc: `Some of the popular ML packages (basic) used in Python are as follows:
              
                  * Numpy.
                  * Pandas.
                  * Matplotlib.
                  * Scikit.
                  
                  Learn each of them; Read from the official documentation of the packages. See the various implementations, look for the WHYs and the WHY NOTs.`,
        url: '',
        time: '30 小时',
      },
    ],
  },
  {
    title: '前端React工程师',
    knowledgePath: [
      {
        title: 'HTML',
        studyOptions: [
          {
            title: 'Free Code Camp basic-html-and-html5',
            type: ['online-course', 'programming-project'],
            language: 'en',
            url:
              'https://www.freecodecamp.org/learn/responsive-web-design/basic-html-and-html5/',
          },
          {
            title: 'freeCodeCamp HTML 基础',
            type: ['online-course', 'programming-project'],
            language: 'cn',
            url:
              'https://learn.freecodecamp.one/responsive-web-design/basic-html-and-html5',
          },
        ],
      },
      {
        title: '学习 CSS',
        studyOptions: [
          {
            title: 'Free Code Camp basic-css',
            type: ['online courses', 'programming projects'],
            url:
              'https://www.freecodecamp.org/learn/responsive-web-design/basic-css/',
          },
        ],
      },
      {
        title: '学习基础 Javascript',
      },
      {
        title: '学习基础 Javascript',
      },
      {
        title: '学习使用Git',
      },
      {
        title: '学习使用 SSH',
      },
      {
        title: '学习 HTTP / HTTPS',
      },
      {
        title: '学习 Linux / Bash',
      },
      {
        title: '学习 GitHub',
      },
      {
        title: '学习 Character Encodings',
      },
    ],
  },
];
