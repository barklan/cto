import { acceptHMRUpdate, defineStore } from 'pinia'

export const useUserStore = defineStore('user', () => {
  /**
   * Current named of the user.
   */
  const name = ref('')
  const id = ref('')
  const project = ref('')
  const token = ref('')

  const projectColor = ref('bg-blue-900')
  const projectPrettyName = ref('')

  /**
   *
   * @param name - new name to set
   */
  function setName(newName: string) {
    name.value = newName
  }

  function pickColor() {
    var baseColors = ['orange', 'yellow', 'green', 'cyan', 'indigo', 'violet', 'fuchsia']
    var randomBaseColor = baseColors[Math.floor(
      Math.random() * baseColors.length)];

    var randomColor = 'bg-' + randomBaseColor + '-900'
    if (randomColor == projectColor.value) {
      randomColor = pickColor()
    }
    return randomColor
  }

  // TODO
  function getProjectPrettyName(projectID) {
    return ""
  }


  function setID(newID: string) {
    id.value = newID
  }
  function setProject(newProject: string) {
    project.value = newProject
    projectColor.value = pickColor()
    projectPrettyName.value = getProjectPrettyName(project.value)
  }
  function setToken(newToken: string) {
    token.value = newToken
  }

  return {
    setName,
    setID,
    setProject,
    setToken,
    name,
    id,
    project,
    projectColor,
    projectPrettyName,
    token,
  }
})

if (import.meta.hot)
  import.meta.hot.accept(acceptHMRUpdate(useUserStore, import.meta.hot))
