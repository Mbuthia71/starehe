import { View, Text, ScrollView, StyleSheet, TouchableOpacity, TextInput } from 'react-native';
import { useColorScheme } from 'react-native';
import { useState } from 'react';
import { Colors } from '@/constants/theme';
import { getJobs, createJob } from '@/lib/api';

export default function JobsScreen() {
  const scheme = useColorScheme();
  const colors = Colors[scheme === 'unspecified' ? 'light' : scheme];
  const [searchQuery, setSearchQuery] = useState('');
  const [showModal, setShowModal] = useState(false);

  return (
    <View style={[styles.container, { backgroundColor: colors.background }]}>
      <ScrollView style={styles.scrollView}>
        <View style={styles.header}>
          <Text style={[styles.title, { color: colors.text }]}>Jobs</Text>
          <Text style={[styles.subtitle, { color: colors.textSecondary }]}>
            Career opportunities for Old Starehians
          </Text>
        </View>

        <View style={styles.searchContainer}>
          <TextInput
            style={[styles.searchInput, { color: colors.text, backgroundColor: colors.card }]}
            placeholder="Search jobs..."
            placeholderTextColor={colors.textSecondary}
            value={searchQuery}
            onChangeText={setSearchQuery}
          />
        </View>

        <View style={styles.content}>
          <View style={styles.headerRow}>
            <Text style={[styles.sectionTitle, { color: colors.text }]}>
              Available Jobs
            </Text>
            <TouchableOpacity
              style={[styles.addButton, { backgroundColor: colors.primary }]}
              onPress={() => setShowModal(true)}
            >
              <Text style={styles.addButtonText}>+ Post Job</Text>
            </TouchableOpacity>
          </View>
          
          {/* Sample job cards */}
          {[
            { title: 'Senior Software Engineer', company: 'Tech Corp', location: 'Nairobi', type: 'Full-time' },
            { title: 'Marketing Manager', company: 'Brand Ltd', location: 'Mombasa', type: 'Full-time' },
            { title: 'Data Analyst', company: 'Data Inc', location: 'Remote', type: 'Contract' },
          ].map((job, index) => (
            <TouchableOpacity
              key={index}
              style={[styles.card, { backgroundColor: colors.card, borderColor: colors.border }]}
            >
              <View style={styles.cardHeader}>
                <Text style={[styles.jobTitle, { color: colors.text }]}>{job.title}</Text>
                <Text style={[styles.jobType, { color: colors.primary }]}>{job.type}</Text>
              </View>
              <Text style={[styles.company, { color: colors.textSecondary }]}>{job.company}</Text>
              <Text style={[styles.location, { color: colors.textSecondary }]}>{job.location}</Text>
            </TouchableOpacity>
          ))}
        </View>
      </ScrollView>
    </View>
  );
}

const styles = StyleSheet.create({
  container: {
    flex: 1,
  },
  scrollView: {
    flex: 1,
  },
  header: {
    padding: 20,
    paddingTop: 60,
  },
  title: {
    fontSize: 32,
    fontWeight: 'bold',
  },
  subtitle: {
    fontSize: 14,
    marginTop: 4,
  },
  searchContainer: {
    marginHorizontal: 20,
    marginBottom: 20,
  },
  searchInput: {
    height: 48,
    borderRadius: 12,
    paddingHorizontal: 16,
    fontSize: 14,
  },
  content: {
    padding: 20,
  },
  headerRow: {
    flexDirection: 'row',
    justifyContent: 'space-between',
    alignItems: 'center',
    marginBottom: 16,
  },
  sectionTitle: {
    fontSize: 18,
    fontWeight: '600',
  },
  addButton: {
    paddingHorizontal: 16,
    paddingVertical: 8,
    borderRadius: 8,
  },
  addButtonText: {
    color: '#fff',
    fontSize: 12,
    fontWeight: '600',
  },
  card: {
    padding: 16,
    borderRadius: 12,
    marginBottom: 12,
    borderWidth: 1,
  },
  cardHeader: {
    flexDirection: 'row',
    justifyContent: 'space-between',
    alignItems: 'flex-start',
    marginBottom: 8,
  },
  jobTitle: {
    fontSize: 16,
    fontWeight: '600',
    flex: 1,
  },
  jobType: {
    fontSize: 12,
    fontWeight: '500',
  },
  company: {
    fontSize: 14,
    marginTop: 4,
  },
  location: {
    fontSize: 12,
    marginTop: 2,
  },
});
