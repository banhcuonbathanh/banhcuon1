import { ReadingTest, Section, Passage, Question, QuestionType } from "./types";

export const readingTest1: ReadingTest = {
  testNumber: 1,
  sections: [
    {
      sectionNumber: 1,
      timeAllowed: 20,
      passages: [
        {
          passageNumber: 1,
          title: "Ants Could Teach Ants",
          content: [
            {
              A: "The ants are tiny and usually nest between rocks in the south coast of England. Transformed into research subjects at the University of Bristol, they raced along a tabletop foraging for food - and then, remarkably, returned to guide others. Time and again, followers trailed behind leaders, darting this way and that along the route, presumably to memorise landmarks. Once a follower got its bearings, it tapped the leader with its antennae, prompting the lesson to literally proceed to the next step. The ants were only looking for food, but the researchers said the careful way the leaders led followers, thereby turning them into leaders in their own right, marked the Temnothorax albipennis ant as the very first example of a non-human animal exhibiting teaching behaviour."
            },
            {
              B: "\"Tandem running is an example of teaching, to our knowledge the first in a non-human animal, that involves bidirectional feedback between teacher and pupil\" remarks Nigel Franks, professor of animal behaviour and ecology, whose paper on the ant educators was published last week in the journal Nature."
            },
            {
              C: "No sooner was the paper published, of course, than another educator questioned it. Marc Hauser, a psychologist and biologist and one of the scientists who came up with the definition of teaching, said it was unclear whether the ants had learned a new skill or merely acquired new information."
            },
            {
              D: "Later, Franks took a further study and found that there were even races between leaders. With the guidance of leaders, ants could find food faster. But the help comes at a cost for the leader, who normally would have reached the food about four times faster if not hampered by a follower. This means the hypothesis that the leaders deliberately slowed down in order to pass the skills on to the followers seems potentially valid. His ideas were advocated by the students who carried out the video project with him."
            },
            {
              E: "Opposing views still arose, however. Hauser noted that mere communication of information is commonplace in the animal world. Consider a species, for example, that uses alarm calls to warn fellow members about the presence. Sounding the alarm can be costly, because the animal may draw the attention of the predator to itself. But it allows others flee to safety. \"Would you call this teaching?\" wrote Hauser. \"The caller incurs a cost. The naive animals gain a benefit and new knowledge that better enables them to learn about the predator's location than if the caller had not called. This happens throughout the animal kingdom, but we don't call it teaching, even though it is clearly transfer of information.\""
            },
            {
              F: "Tim Caro, a zoologist, presented two cases of animal communication. He found that cheetah mothers that take their cubs along on hunts gradually allow their cubs to do more of the hunting —going, for example, from killing a gazelle and allowing young cubs to eat merely tripping the gazelle and letting the cubs finish it off. At one level, such behaviour might be called teaching — except the mother was not really teaching the cubs to hunt but merely facilitating various stages of learning. In another instance, birds watching other birds using a stick to locate food such as insects and so on, are observed to do the same thing themselves while finding food later."
            },
            {
              G: "Psychologists study animal behaviour in part to understand the evolutionary roots of human behaviour, Hauser said. The challenge in understanding whether other animals truly teach one another, he added, is that human teaching involves a \"theory of mind\" teachers are aware that students don't know something. He questioned whether Franks' leader ants really knew that the follower ants were ignorant. Could they simply have been following an instinctive rule to proceed when the followers tapped them on the legs or abdomen? And did leaders that led the way to food 一 only to find that it had been removed by the experimenter - incur the wrath of followers? That, Hauser said, would suggest that the follower ant actually knew the leader was more knowledgeable and not merely following an instinctive routine itself."
            },
            {
              H: "The controversy went on, and for a good reason. The occurrence of teaching in ants, if proven to be true, indicates that teaching can evolve in animals with tiny brains. It is probably the value of information in social animals that determines when teaching will evolve, rather than the constraints of brain size."
            },
            {
              I: "Bennett Galef Jr., a psychologist who studies animal behaviour and social learning at McMaster University in Canada，maintained that ants were unlikely to have a \"theory of mind\" 一 meaning that leaders and followers may well have been following instinctive routines that were not based on an understanding of what was happening in another ant's brain. He warned that scientists may be barking up the wrong tree when they look not only for examples of humanlike behaviour among other animals but humanlike thinking that underlies such behaviour. Animals may behave in ways similar to humans without a similar cognitive system, he said, so the behaviour is not necessarily a good guide into how humans came to think the way they do."
            }
          ],
          questions: [
            {
              questionNumber: 1,
              type: QuestionType.Matching,
              content: "Animals could use objects to locate food.",
              options: [
                "A Nigel Franks",
                "B Marc Hauser",
                "C Tim Caro",
                "D Bennet Galef Jr"
              ],
              correctAnswer: "C"
            },
            {
              questionNumber: 2,
              type: QuestionType.Matching,
              content: "Ants show two-way, interactive teaching behaviours.",
              options: [
                "A Nigel Franks",
                "B Marc Hauser",
                "C Tim Caro",
                "D Bennet Galef Jr"
              ],
              correctAnswer: "A"
            },
            {
              questionNumber: 3,
              type: QuestionType.Matching,
              content: "It is risky to say ants can teach other ants like human beings do.",
              options: [
                "A Nigel Franks",
                "B Marc Hauser",
                "C Tim Caro",
                "D Bennet Galef Jr"
              ],
              correctAnswer: "D"
            },
            {
              questionNumber: 4,
              type: QuestionType.Matching,
              content: "Ant leadership makes finding food faster.",
              options: [
                "A Nigel Franks",
                "B Marc Hauser",
                "C Tim Caro",
                "D Bennet Galef Jr"
              ],
              correctAnswer: "A"
            },
            {
              questionNumber: 5,
              type: QuestionType.Matching,
              content: "Communication between ants is not entirely teaching.",
              options: [
                "A Nigel Franks",
                "B Marc Hauser",
                "C Tim Caro",
                "D Bennet Galef Jr"
              ],
              correctAnswer: "B"
            },
            {
              questionNumber: 6,
              type: QuestionType.MultipleChoice,
              content: "Which FOUR of the following behaviours of animals are mentioned in the passage?",
              options: [
                "A touch each other with antenna",
                "B alert others when there is danger",
                "C escape from predators",
                "D protect the young",
                "E hunt food for the young",
                "F fight with each other",
                "G use tools like twigs",
                "H feed on a variety of foods"
              ],
              correctAnswer: ["A", "B", "C", "G"]
            },
            {
              questionNumber: 10,
              type: QuestionType.TrueFalseNotGiven,
              content: "Ants' tandem running involves only one-way communication.",
              correctAnswer: "FALSE"
            },
            {
              questionNumber: 11,
              type: QuestionType.TrueFalseNotGiven,
              content: "Franks's theory got many supporters immediately after publicity.",
              correctAnswer: "FALSE"
            },
            {
              questionNumber: 12,
              type: QuestionType.TrueFalseNotGiven,
              content: "Ants' teaching behaviour is the same as that of human.",
              correctAnswer: "FALSE"
            },
            {
              questionNumber: 13,
              type: QuestionType.TrueFalseNotGiven,
              content: "Cheetah share hunting gains to younger ones",
              correctAnswer: "TRUE"
            }
          ]
        }
      ]
    }
  ]
};